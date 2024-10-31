package ratelimiter

type ServiceMock struct {
	AllowRequestFunc    func(identifier string, reqType string) (bool, error)
	GetRequestCountFunc func(identifier string) (int, error)
	GetTTLFunc          func(identifier string) (int, error)
	AlreadyExistsFunc   func(identifier string) (bool, error)
}

func (s ServiceMock) AllowRequest(identifier string, reqType string) (bool, error) {
	if s.AllowRequestFunc != nil {
		return s.AllowRequestFunc(identifier, reqType)
	}
	panic("AllowRequestFunc is not defined for ServiceMock")
}

func (s ServiceMock) GetRequestCount(identifier string) (int, error) {
	if s.GetRequestCountFunc != nil {
		return s.GetRequestCountFunc(identifier)
	}
	panic("GetRequestCountFunc is not defined for ServiceMock")
}

func (s ServiceMock) GetTTL(identifier string) (int, error) {
	if s.GetTTLFunc != nil {
		return s.GetTTLFunc(identifier)
	}
	panic("GetTTLFunc is not defined for ServiceMock")
}

func (s ServiceMock) AlreadyExists(identifier string) (bool, error) {
	if s.AlreadyExistsFunc != nil {
		return s.AlreadyExistsFunc(identifier)
	}
	panic("AlreadyExistsFunc is not defined for ServiceMock")
}
