package main

import (
	"sync"
	"time"
)

type aircraftCache struct {
	mu   sync.Mutex
	data map[string]aircraftData
	ages map[string]time.Time
	ttl  time.Duration
}

// newAircraftCache returns an [aircraftCache] pointer
// and starts a background ttl monitor.
func newAircraftCache(ttl time.Duration) *aircraftCache {
	cache := &aircraftCache{
		data: make(map[string]aircraftData),
		ages: make(map[string]time.Time),
		ttl:  ttl,
	}

	// Monitor once a second
	go cache.ttlMonitor(time.Second)

	return cache
}

func (c *aircraftCache) set(callsign string, data aircraftData) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[callsign] = data
	c.ages[callsign] = time.Now()
}

func (c *aircraftCache) get(callsign string) (aircraftData, bool) {
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
			if age.Add(c.ttl).After(time.Now()) {
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
		}
		c.mu.Unlock()
	}
}
