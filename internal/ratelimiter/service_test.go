package ratelimiter_test

import (
	"github.com/ThamirisMonteiro/rate-limiter/internal/infra/database"
	"github.com/ThamirisMonteiro/rate-limiter/internal/ratelimiter"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type ServiceTestSuite struct {
	suite.Suite
	repository database.Repository
	service    *ratelimiter.Service
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.repository = database.NewRedisRepository("localhost:6379")
	suite.NoError(suite.repository.FlushDB())
	err := os.Setenv("REQ_LIMIT", "2")
	suite.NoError(err)
	suite.service = ratelimiter.NewService(suite.repository)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) TestService_SetLimit() {
	err := suite.service.SetLimit(10, time.Second)
	suite.NoError(err)

	ttl, err := suite.repository.TTLKey(ratelimiter.SettingsKey)
	suite.NoError(err)
	suite.Equal(time.Second, ttl)
}

func (suite *ServiceTestSuite) TestService_GetRequestCount() {
	err := suite.repository.SetKey("requests:test", "1", 0)
	suite.NoError(err)

	count, err := suite.service.GetRequestCount("test")
	suite.NoError(err)
	suite.Equal(1, count)
}

func (suite *ServiceTestSuite) TestService_AllowRequest() {
	identifier := "test"

	err := suite.repository.IncrCounter("requests:" + identifier)
	suite.Require().NoError(err)

	allowed, err := suite.service.AllowRequest(identifier)
	suite.NoError(err)
	suite.True(allowed)

	count, err := suite.service.GetRequestCount(identifier)
	suite.NoError(err)
	suite.Equal(2, count)

	allowed, err = suite.service.AllowRequest(identifier)
	suite.NoError(err)
	suite.False(allowed)

	count, err = suite.service.GetRequestCount(identifier)
	suite.NoError(err)
	suite.Equal(2, count)
}
