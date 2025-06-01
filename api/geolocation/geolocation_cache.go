package geolocation

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

var cache = make(map[string]CacheEntry)
var mu sync.Mutex

func SetToCache(key string, data interface{}) {
	mu.Lock()
	defer mu.Unlock()
	cache[key] = CacheEntry{Data: data, ExpiresAt: time.Now().Add(24 * time.Hour)}
}

func GetFromCache(key string) (interface{}, bool) {
	mu.Lock()
	defer mu.Unlock()
	entry, exists := cache[key]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Data, true
}
