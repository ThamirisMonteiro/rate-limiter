package ratelimiter

import "time"

type RateLimiter interface {
	AllowRequest(identifier string) (bool, error)
	GetRequestCount(identifier string) (int, error)
	SetLimit(key string, limit int, ttl time.Duration) error
}
