package main

import (
	"sync"
	"time"
)

type aircraftCache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
	ttl   time.Duration
}

type cacheItem struct {
	data        AircraftData
	lastUpdated time.Time
}

type AircraftData struct {
	Callsign          string    `json:"callsign"`
	AircraftType      string    `json:"aircraft_type"`
	PilotName         string    `json:"pilot_name"`
	GroupName         string    `json:"group_name"`
	MsfsServer        string    `json:"msfs_server"`
	TransponderCode   string    `json:"transponder_code"`
	Latitude          string    `json:"latitude"`
	Longitude         string    `json:"longitude"`
	Altitude          string    `json:"altitude"`
	Heading           string    `json:"heading"`
	Airspeed          string    `json:"airspeed"`
	Groundspeed       string    `json:"groundspeed"`
	TouchdownVelocity string    `json:"touchdown_velocity"`
	Notes             string    `json:"notes"`
	Version           string    `json:"version"`
	Created           time.Time `json:"-"`
	Modified          time.Time `json:"-"`
}

// newAircraftCache returns an [aircraftCache] pointer
// and starts a background ttl monitor.
func newAircraftCache(ttl time.Duration) *aircraftCache {
	cache := &aircraftCache{
		items: make(map[string]cacheItem),
		ttl:   ttl,
	}

	// Monitor once a second
	go cache.startTTLMonitor(time.Second)

	return cache
}

func (c *aircraftCache) set(callsign string, data AircraftData) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := cacheItem{
		data:        data,
		lastUpdated: time.Now(),
	}

	c.items[callsign] = item
}

func (c *aircraftCache) get(callsign string) (item AircraftData, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, ok := c.items[callsign]

	return data.data, ok
}

func (c *aircraftCache) startTTLMonitor(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for range t.C {
		now := time.Now()

		c.mu.Lock()
		for callsign, item := range c.items {
			if item.lastUpdated.Add(c.ttl).Before(now) {
				// Delete the expired item
				delete(c.items, callsign)
			}
		}
		c.mu.Unlock()
	}
}
