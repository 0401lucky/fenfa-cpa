package middleware

import (
	"cpa-distribution/common/utils"
	"cpa-distribution/model"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type slidingWindow struct {
	timestamps []time.Time
	mu         sync.Mutex
}

var (
	rateLimitWindows = make(map[uint]*slidingWindow)
	rlMutex          sync.RWMutex
)

func getWindow(tokenID uint) *slidingWindow {
	rlMutex.RLock()
	w, exists := rateLimitWindows[tokenID]
	rlMutex.RUnlock()
	if exists {
		return w
	}

	rlMutex.Lock()
	defer rlMutex.Unlock()
	if w, exists = rateLimitWindows[tokenID]; exists {
		return w
	}
	w = &slidingWindow{}
	rateLimitWindows[tokenID] = w
	return w
}

func RateLimit() gin.HandlerFunc {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			cleanupWindows()
		}
	}()

	return func(c *gin.Context) {
		tokenRaw, exists := c.Get("token")
		if !exists {
			c.Next()
			return
		}
		token := tokenRaw.(*model.Token)
		if token.RateLimitRPM <= 0 {
			c.Next()
			return
		}

		w := getWindow(token.ID)
		w.mu.Lock()

		now := time.Now()
		windowStart := now.Add(-time.Minute)

		valid := make([]time.Time, 0, len(w.timestamps))
		for _, ts := range w.timestamps {
			if ts.After(windowStart) {
				valid = append(valid, ts)
			}
		}
		w.timestamps = valid

		if len(w.timestamps) >= token.RateLimitRPM {
			w.mu.Unlock()
			utils.SendOpenAIError(c, http.StatusTooManyRequests, "rate_limit_exceeded", "Rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}

		w.timestamps = append(w.timestamps, now)
		w.mu.Unlock()
		c.Next()
	}
}

func cleanupWindows() {
	rlMutex.Lock()
	defer rlMutex.Unlock()
	now := time.Now()
	for id, w := range rateLimitWindows {
		w.mu.Lock()
		if len(w.timestamps) == 0 || w.timestamps[len(w.timestamps)-1].Before(now.Add(-2*time.Minute)) {
			delete(rateLimitWindows, id)
		}
		w.mu.Unlock()
	}
}
