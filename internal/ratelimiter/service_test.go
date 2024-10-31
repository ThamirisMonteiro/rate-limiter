package ratelimiter_test

import (
	"fmt"
	"github.com/ThamirisMonteiro/rate-limiter/internal/infra/database"
	"github.com/ThamirisMonteiro/rate-limiter/internal/ratelimiter"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ServiceTestSuite struct {
	suite.Suite
}

func (suite *ServiceTestSuite) SetupTest() {
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) Test_AlreadyExists_WhenGetKeyFails_ShouldReturnError() {
	reqLimit := 1
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		GetKeyFunc: func(key string) (string, error) {
			return "", fmt.Errorf("error getting key")
		},
		SetKeyFunc: func(key string, value string, expiration time.Duration) error {
			return nil
		},
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	exists, err := service.AlreadyExists("key")

	suite.Error(err)
	suite.False(exists)
	suite.Contains(err.Error(), "error getting key")
}

func (suite *ServiceTestSuite) Test_AlreadyExists_WhenKeyDoesNotExist_ShouldReturnFalse() {
	reqLimit := 1
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		GetKeyFunc: func(key string) (string, error) {
			return "", fmt.Errorf("redis: nil")
		},
		SetKeyFunc: func(key string, value string, expiration time.Duration) error {
			return nil
		},
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	exists, err := service.AlreadyExists("key")

	suite.NoError(err)
	suite.False(exists)
}

func (suite *ServiceTestSuite) Test_AlreadyExists_WhenKeyExists_ShouldReturnTrue() {
	reqLimit := 1
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		GetKeyFunc: func(key string) (string, error) {
			return "key", nil
		},
		SetKeyFunc: func(key string, value string, expiration time.Duration) error {
			return nil
		},
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	exists, err := service.AlreadyExists("key")

	suite.NoError(err)
	suite.True(exists)
}

func (suite *ServiceTestSuite) Test_GetTTL_IfGetTTLFails_ShouldReturnError() {
	reqLimit := 1
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		SetKeyFunc: func(key string, value string, expiration time.Duration) error {
			return nil
		},
		TTLKeyFunc: func(key string) (time.Duration, error) {
			return time.Duration(0), fmt.Errorf("error getting ttl")
		},
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	_, err = service.GetTTL("key")

	suite.Error(err)
	suite.Contains(err.Error(), "error getting ttl")
}

func (suite *ServiceTestSuite) Test_GetTTL() {
	reqLimit := 1
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		SetKeyFunc: func(key string, value string, expiration time.Duration) error {
			return nil
		},
		TTLKeyFunc: func(key string) (time.Duration, error) {
			return time.Duration(1), nil
		},
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	ttl, err := service.GetTTL("key")

	suite.NoError(err)
	suite.Equal(1, ttl)
}

func (suite *ServiceTestSuite) Test_GetRequestCount_WhenGetKeyFails_ShouldReturnError() {
	reqLimit := 1
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		SetKeyFunc: func(key string, value string, expiration time.Duration) error {
			return nil
		},
		TTLKeyFunc: func(key string) (time.Duration, error) {
			return time.Duration(1), nil
		},
		GetKeyFunc: func(key string) (string, error) {
			return "", fmt.Errorf("error getting key")
		},
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	_, err = service.GetRequestCount("key")

	suite.Error(err)
	suite.Contains(err.Error(), "error getting key")
}

func (suite *ServiceTestSuite) Test_GetRequestCount_WhenAtoiFails_ShouldReturnError() {
	reqLimit := 1
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		SetKeyFunc: func(key string, value string, expiration time.Duration) error {
			return nil
		},
		TTLKeyFunc: func(key string) (time.Duration, error) {
			return time.Duration(1), nil
		},
		GetKeyFunc: func(key string) (string, error) {
			return "key", nil
		},
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	_, err = service.GetRequestCount("key")

	suite.Error(err)
	suite.Contains(err.Error(), "invalid syntax")
}

func (suite *ServiceTestSuite) Test_GetRequestCount_WhenKeyDoesNotExist_ShouldDoNothing() {
	reqLimit := 1
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		SetKeyFunc: func(key string, value string, expiration time.Duration) error {
			return nil
		},
		TTLKeyFunc: func(key string) (time.Duration, error) {
			return time.Duration(1), nil
		},
		GetKeyFunc: func(key string) (string, error) {
			return "", fmt.Errorf("redis: nil")
		},
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	requestCount, err := service.GetRequestCount("key")

	suite.NoError(err)
	suite.Equal(0, requestCount)
}

func (suite *ServiceTestSuite) Test_GetRequestCount() {
	reqLimit := 1
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		SetKeyFunc: func(key string, value string, expiration time.Duration) error {
			return nil
		},
		TTLKeyFunc: func(key string) (time.Duration, error) {
			return time.Duration(1), nil
		},
		GetKeyFunc: func(key string) (string, error) {
			return "1", nil
		},
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	requestCount, err := service.GetRequestCount("key")

	suite.NoError(err)
	suite.Equal(1, requestCount)
}

func (suite *ServiceTestSuite) Test_NewService() {
	reqLimit := 1
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		SetKeyFunc: func(key string, value string, expiration time.Duration) error {
			return nil
		},
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	suite.NotNil(service)
}

func (suite *ServiceTestSuite) Test_NewService_WhenSetKeyFails_ShouldReturnError() {
	reqLimit := 1
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		SetKeyFunc: func(key string, value string, expiration time.Duration) error {
			return fmt.Errorf("error setting key")
		},
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.Error(err)

	suite.Nil(service)
	suite.Contains(err.Error(), "error setting key")
}

func (suite *ServiceTestSuite) Test_AllowRequest_WhenRequestIsAllowed_ShouldReturnTrue() {
	reqLimit := 2
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		SetKeyFunc: func(key string, value string, expiration time.Duration) error {
			return nil
		},
		TTLKeyFunc: func(key string) (time.Duration, error) {
			return time.Duration(1), nil
		},
		GetKeyFunc: func(key string) (string, error) {
			return "1", nil
		},
		IncrCounterFunc: func(counter string) error { return nil },
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	isAllowed, err := service.AllowRequest("key", "ip")

	suite.NoError(err)
	suite.True(isAllowed)
}

func (suite *ServiceTestSuite) Test_AllowRequest_WhenAlreadyExistsFails_ShouldReturnError() {
	reqLimit := 2
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		SetKeyFunc: func(key string, value string, expiration time.Duration) error {
			return nil
		},
		TTLKeyFunc: func(key string) (time.Duration, error) {
			return time.Duration(1), nil
		},
		GetKeyFunc: func(key string) (string, error) {
			return "", fmt.Errorf("unable to get key")
		},
		IncrCounterFunc: func(counter string) error { return nil },
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	isAllowed, err := service.AllowRequest("key", "ip")

	suite.Error(err)
	suite.False(isAllowed)
	suite.Contains(err.Error(), "unable to get key")
}

func (suite *ServiceTestSuite) Test_AllowRequest_WhenAlreadyExistsIsFalseAndIsIP_ShouldReturnTrue() {
	reqLimit := 2
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		SetKeyFunc: func(key string, value string, expiration time.Duration) error { return nil },
		TTLKeyFunc: func(key string) (time.Duration, error) {
			return time.Duration(1), nil
		},
		GetKeyFunc:      func(key string) (string, error) { return "", fmt.Errorf("redis: nil") },
		IncrCounterFunc: func(counter string) error { return nil },
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	isAllowed, err := service.AllowRequest("key", "ip")

	suite.NoError(err)
	suite.True(isAllowed)
}

func (suite *ServiceTestSuite) Test_AllowRequest_WhenAlreadyExistsIsFalseAndIsToken_ShouldReturnTrue() {
	reqLimit := 2
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(1)

	repo := database.RepositoryMock{
		SetKeyFunc: func(key string, value string, expiration time.Duration) error { return nil },
		TTLKeyFunc: func(key string) (time.Duration, error) {
			return time.Duration(1), nil
		},
		GetKeyFunc:      func(key string) (string, error) { return "", fmt.Errorf("redis: nil") },
		IncrCounterFunc: func(counter string) error { return nil },
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	isAllowed, err := service.AllowRequest("key", "token")

	suite.NoError(err)
	suite.True(isAllowed)
}

func (suite *ServiceTestSuite) Test_AllowRequest_WhenAlreadyExistsIsFalseAndTokenTimeIsBigger_ShouldReturnTrue() {
	reqLimit := 2
	blockTimeIP := time.Duration(1)
	blockTimeToken := time.Duration(10)

	repo := database.RepositoryMock{
		SetKeyFunc: func(key string, value string, expiration time.Duration) error { return nil },
		TTLKeyFunc: func(key string) (time.Duration, error) {
			return time.Duration(1), nil
		},
		GetKeyFunc:      func(key string) (string, error) { return "", fmt.Errorf("redis: nil") },
		IncrCounterFunc: func(counter string) error { return nil },
	}

	service, err := ratelimiter.NewService(repo, reqLimit, blockTimeIP, blockTimeToken)
	suite.NoError(err)

	isAllowed, err := service.AllowRequest("key", "token")

	suite.NoError(err)
	suite.True(isAllowed)
}
