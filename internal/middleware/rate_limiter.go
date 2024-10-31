package middleware

import (
	"github.com/ThamirisMonteiro/rate-limiter/internal/ratelimiter"
	"net"
	"net/http"
)

func RateLimiterMiddleware(service *ratelimiter.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			identifier := extractIdentifier(r)

			allowed := checkRequestLimit(service, identifier, w)
			if !allowed {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractIdentifier(r *http.Request) string {
	token := r.Header.Get("API_KEY")
	if token != "" {
		return token
	}

	ip := r.RemoteAddr

	host, _, err := net.SplitHostPort(ip)
	if err != nil {
		return ip
	}

	return host
}

func checkRequestLimit(service *ratelimiter.Service, identifier string, w http.ResponseWriter) bool {
	allowed, err := service.AllowRequest(identifier)
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
