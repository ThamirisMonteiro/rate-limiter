package database

import "time"

type Repository interface {
	SetKey(key string, value string, expiration time.Duration) error
	GetKey(key string) (string, error)
	IncrCounter(counter string) error
	ExpireKey(key string, expiration time.Duration) error
	TTLKey(key string) (time.Duration, error)
}
