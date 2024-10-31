package ratelimiter

type RateLimiter interface {
	AllowRequest(identifier string, reqType string) (bool, error)
	GetRequestCount(identifier string) (int, error)
	GetTTL(identifier string) (int, error)
	AlreadyExists(identifier string) (bool, error)
}
