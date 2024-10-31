package middleware_test

import (
	"github.com/ThamirisMonteiro/rate-limiter/internal/middleware"
	"github.com/ThamirisMonteiro/rate-limiter/internal/ratelimiter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RateLimiterTestSuite struct {
	suite.Suite
}

func (suite *RateLimiterTestSuite) SetupTest() {
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterTestSuite))
}

func (suite *RateLimiterTestSuite) Test_RateLimiterMiddleware() {
	service := &ratelimiter.ServiceMock{
		AllowRequestFunc: func(identifier string, reqType string) (bool, error) {
			return true, nil
		},
	}
	rlmiddleware := middleware.RateLimiterMiddleware(service)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Request successful"))
	})

	rlmiddleware(handler).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		suite.T().Errorf("Expected status %d but got %d", http.StatusOK, w.Code)
	}
}

func (suite *RateLimiterTestSuite) Test_RateLimiterMiddleware_TooManyRequests() {
	service := &ratelimiter.ServiceMock{
		AllowRequestFunc: func(identifier string, reqType string) (bool, error) {
			return false, nil
		},
	}
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	})

	rlmiddleware := middleware.RateLimiterMiddleware(service)(nextHandler)

	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("API_KEY", "test-api-key")

	rr := httptest.NewRecorder()

	rlmiddleware.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusTooManyRequests, rr.Code)
	assert.Equal(suite.T(), "You have reached the maximum number of requests or actions allowed within a certain time frame", rr.Body.String())
}
