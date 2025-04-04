package rate_limiter

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"

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

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}
