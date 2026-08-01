package main

import (
	"sync"
	"time"
)

type cache[T any] struct {
	mu    sync.RWMutex
	items map[string]cacheItem[T]
	ttl   time.Duration
}

type cacheItem[T any] struct {
	data        T
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

// newCache returns a [cache] pointer
// and starts a background TTL monitor.
func newCache[T any](ttl time.Duration) *cache[T] {
	cache := &cache[T]{
		items: make(map[string]cacheItem[T]),
		ttl:   ttl,
	}

	// Monitor once a second
	go cache.startTTLMonitor(time.Second)

	return cache
}

func (c *cache[T]) set(key string, data T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem[T]{
		data:        data,
		lastUpdated: time.Now().UTC(),
	}
}

func (c *cache[T]) get(key string) (item T, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, ok := c.items[key]

	return data.data, ok
}

func (c *cache[T]) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Reinitialise the map to free memory
	c.items = make(map[string]cacheItem[T])
}

func (c *cache[T]) len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.items)
}

func (c *cache[T]) startTTLMonitor(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for range t.C {
		now := time.Now().UTC()

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
