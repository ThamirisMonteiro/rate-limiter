package database

import "time"

type RepositoryMock struct {
	SetKeyFunc      func(key string, value string, expiration time.Duration) error
	GetKeyFunc      func(key string) (string, error)
	IncrCounterFunc func(counter string) error
	ExpireKeyFunc   func(key string, expiration time.Duration) error
	TTLKeyFunc      func(key string) (time.Duration, error)
	CloseFunc       func() error
	FlushDBFunc     func() error
}

func (r RepositoryMock) SetKey(key string, value string, expiration time.Duration) error {
	if r.SetKeyFunc != nil {
		return r.SetKeyFunc(key, value, expiration)
	}
	panic("SetKeyFunc is not defined for RepositoryMock")
}

func (r RepositoryMock) GetKey(key string) (string, error) {
	if r.GetKeyFunc != nil {
		return r.GetKeyFunc(key)
	}
	panic("GetKeyFunc is not defined for RepositoryMock")
}

func (r RepositoryMock) IncrCounter(counter string) error {
	if r.IncrCounterFunc != nil {
		return r.IncrCounterFunc(counter)
	}
	panic("IncrCounterFunc is not defined for RepositoryMock")
}

func (r RepositoryMock) ExpireKey(key string, expiration time.Duration) error {
	if r.ExpireKeyFunc != nil {
		return r.ExpireKeyFunc(key, expiration)
	}
	panic("ExpireKeyFunc is not defined for RepositoryMock")
}

func (r RepositoryMock) TTLKey(key string) (time.Duration, error) {
	if r.TTLKeyFunc != nil {
		return r.TTLKeyFunc(key)
	}
	panic("TTLKeyFunc is not defined for RepositoryMock")
}

func (r RepositoryMock) Close() error {
	if r.CloseFunc != nil {
		return r.CloseFunc()
	}
	panic("CloseFunc is not defined for RepositoryMock")
}

func (r RepositoryMock) FlushDB() error {
	if r.FlushDBFunc != nil {
		return r.FlushDBFunc()
	}
	panic("FlushDBFunc is not defined for RepositoryMock")
}
