package middleware

import (
	"github.com/ThamirisMonteiro/rate-limiter/internal/ratelimiter"
	"net"
	"net/http"
)

func RateLimiterMiddleware(service ratelimiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identifier, reqType := extractIdentifier(r)

			allowed := checkRequestLimit(service, identifier, w, reqType)
			if !allowed {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractIdentifier(r *http.Request) (string, string) {
	var reqType string
	token := r.Header.Get("API_KEY")
	if token != "" {
		reqType = "token"
		return token, reqType
	}

	ip := r.RemoteAddr
	reqType = "ip"

	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		return ip, reqType
	}

	return host, reqType
}

func checkRequestLimit(service ratelimiter.RateLimiter, identifier string, w http.ResponseWriter, reqType string) bool {
	allowed, err := service.AllowRequest(identifier, reqType)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return false
	}

	if !allowed {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("You have reached the maximum number of requests or actions allowed within a certain time frame"))
		return false
	}

	return true
}
