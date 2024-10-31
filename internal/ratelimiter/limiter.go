package ratelimiter

type RateLimiter interface {
	AllowRequest(identifier string) (bool, error)
	GetRequestCount(identifier string) (int, error)
	GetTTL(identifier string) (int, error)
}
