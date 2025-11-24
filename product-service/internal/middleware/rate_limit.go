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

var (
	visitors = make(map[string]*bucket)
	mu       sync.Mutex
)

func getBucket(ip string) *bucket {
	mu.Lock()
	defer mu.Unlock()

	b, exists := visitors[ip]
	if !exists {
		b = &bucket{tokens: 10, last: time.Now()}
		visitors[ip] = b
	}
	return b
}

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		b := getBucket(ip)

		// Refill tokens per second
		if time.Since(b.last) >= time.Second {
			b.tokens = 10
			b.last = time.Now()
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
