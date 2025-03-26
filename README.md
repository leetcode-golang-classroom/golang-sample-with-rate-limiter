# golang-sample-with-rate-limiter

This repository is to demo how to use golang standard library to implement rate limiter

## logic

1. install package
```shell
go get golang.org/x/time/rate
```

2. setup rate limiter
```golang
package rate_limiter

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"golang.org/x/time/rate"
)

type RateLimiter struct{}

// getIP - 取得 request IP
func (r *RateLimiter) getIP(req *http.Request) string {
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		log.Printf("Error parsing IP: %v", err)
		return ""
	}
	return host
}

// RateLimiterMiddleware - 建立 ratelimiter middleware
func (r *RateLimiter) RateLimiterMiddleware(next http.Handler, limit rate.Limit, burst int) http.Handler {
	ipLimiterMap := make(map[string]*rate.Limiter)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Fetch IP
		ip := r.getIP(req)
		// Create limiter if not present for IP
		limiter, exists := ipLimiterMap[ip]
		if !exists {
			limiter = rate.NewLimiter(limit, burst)
			ipLimiterMap[ip] = limiter
		}
		// return error if the limit has been reached
		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": "Too many requests"})
			return
		}
		next.ServeHTTP(w, req)
	})
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}
```

## Fix Current Access with sync.Mutex

```golang
// RateLimiterMiddleware - 建立 ratelimiter middleware
func (r *RateLimiter) RateLimiterMiddleware(next http.Handler, limit rate.Limit, burst int) http.Handler {
	ipLimiterMap := make(map[string]*rate.Limiter)
	var mu sync.Mutex
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Fetch IP
		ip := r.getIP(req)
		// Create limiter if not present for IP
		mu.Lock()
		limiter, exists := ipLimiterMap[ip]
		if !exists {
			limiter = rate.NewLimiter(limit, burst)
			ipLimiterMap[ip] = limiter
		}
		mu.Unlock()
		// return error if the limit has been reached
		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": "Too many requests"})
			return
		}
		next.ServeHTTP(w, req)
	})
}
```


## handle concurrency problem with sync.Map

```golang
var ipLimiterMap sync.Map

// RateLimiterMiddleware - 建立 ratelimiter middleware
func (r *RateLimiter) RateLimiterMiddleware(next http.Handler, limit rate.Limit, burst int) http.Handler {

	// var mu sync.Mutex
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Fetch IP
		ip := r.getIP(req)
		// Create limiter if not present for IP
		limiterAny, _ := ipLimiterMap.LoadOrStore(ip, rate.NewLimiter(limit, burst))
		limiter := limiterAny.(*rate.Limiter)
		// return error if the limit has been reached
		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": "Too many requests"})
			return
		}
		next.ServeHTTP(w, req)
	})
}
```

