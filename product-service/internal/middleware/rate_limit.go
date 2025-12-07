package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type bucket struct {
	tokens int
	last   time.Time
}

type rateLimiter struct {
	visitors map[string]*bucket
	mu       sync.RWMutex
}

var limiter = &rateLimiter{
	visitors: make(map[string]*bucket),
}

func init() {
	// Cleanup old entries every 5 minutes to prevent memory leak
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			limiter.cleanup()
		}
	}()
}

// cleanup removes buckets that haven't been used in the last 10 minutes
func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-10 * time.Minute)
	for ip, b := range rl.visitors {
		if b.last.Before(cutoff) {
			delete(rl.visitors, ip)
		}
	}
}

// getBucket retrieves or creates a bucket for the given IP
func (rl *rateLimiter) getBucket(ip string) *bucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, exists := rl.visitors[ip]
	if !exists {
		b = &bucket{tokens: 10, last: time.Now()}
		rl.visitors[ip] = b
	}
	return b
}

// RateLimit middleware limits requests to 10 per second per IP
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if ip == "" {
			ip = r.RemoteAddr // fallback if SplitHostPort fails
		}

		b := limiter.getBucket(ip)

		// Refill tokens (simple token bucket: 10 tokens per second)
		now := time.Now()
		if now.Sub(b.last) >= time.Second {
			b.tokens = 10
			b.last = now
		}

		if b.tokens <= 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error":"rate limit exceeded"}`))
			return
		}

		b.tokens--
		next.ServeHTTP(w, r)
	})
}
