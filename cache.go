package main

import (
	"sync"
	"time"
)

type aircraftCache struct {
	mu   sync.Mutex
	data map[string]AircraftData
	ages map[string]time.Time
	ttl  time.Duration
}

type AircraftData struct {
	UserPin           string    `json:"-"`
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
	Timestamp         time.Time `json:"-"`
	Created           time.Time `json:"-"`
	Modified          time.Time `json:"-"`
}

// newAircraftCache returns an [aircraftCache] pointer
// and starts a background ttl monitor.
func newAircraftCache(ttl time.Duration) *aircraftCache {
	cache := &aircraftCache{
		data: make(map[string]AircraftData),
		ages: make(map[string]time.Time),
		ttl:  ttl,
	}

	// Monitor once a second
	go cache.ttlMonitor(time.Second)

	return cache
}

func (c *aircraftCache) set(callsign string, data AircraftData) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[callsign] = data
	c.ages[callsign] = time.Now()
}

func (c *aircraftCache) get(callsign string) (AircraftData, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, ok := c.data[callsign]

	return data, ok
}

func (c *aircraftCache) ttlMonitor(interval time.Duration) {
	for {
		time.Sleep(interval)

		// Slice to store expired callsigns
		expired := []string{}

		for callsign, age := range c.ages {
			if age.Add(c.ttl).Before(time.Now()) {
				// The callsign is expired
				expired = append(expired, callsign)
			}
		}

		if len(expired) == 0 {
			// There are no expired callsigns
			continue
		}

		// Delete the expired data
		c.mu.Lock()
		for _, callsign := range expired {
			delete(c.data, callsign)
			delete(c.ages, callsign)
		}
		c.mu.Unlock()
	}
}
