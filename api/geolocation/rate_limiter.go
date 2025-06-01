package geolocation

import (
	"sync"
	"time"
)

var lastRequest = make(map[string]time.Time)
var rlMu sync.Mutex

func AllowRequest(ip string) bool {
	rlMu.Lock()
	defer rlMu.Unlock()

	lastTime, exists := lastRequest[ip]
	if !exists || time.Since(lastTime) > 5*time.Second {
		lastRequest[ip] = time.Now()
		return true
	}
	return false
}
