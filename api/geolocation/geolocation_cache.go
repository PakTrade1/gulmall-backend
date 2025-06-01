package geolocation

import (
	"log"
	"sync"
	"time"
)

type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

var (
	cache = make(map[string]CacheEntry)
	mu    sync.Mutex
)

// Set cache entry with 24hr expiry
func SetToCache(key string, data interface{}) {
	mu.Lock()
	defer mu.Unlock()
	cache[key] = CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
}

// Get from cache and ensure country and country_code exist
func GetFromCache(key string) (interface{}, bool) {
	mu.Lock()
	defer mu.Unlock()

	entry, exists := cache[key]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	// Ensure required fields exist
	if geoData, ok := entry.Data.(map[string]interface{}); ok {
		country, cOk := geoData["country"]
		code, ccOk := geoData["country_code"]

		if !cOk || !ccOk || country == "" || code == "" {
			return nil, false
		}
	}

	return entry.Data, true
}

// StartCleanupRoutine starts a background cleanup every interval
func StartCleanupRoutine(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			cleanupExpiredEntries()
		}
	}()
}

func cleanupExpiredEntries() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	removed := 0

	for k, v := range cache {
		if now.After(v.ExpiresAt) {
			delete(cache, k)
			removed++
		}
	}

	if removed > 0 {
		log.Printf("[Cache Cleanup] Removed %d expired entries\n", removed)
	}
}
