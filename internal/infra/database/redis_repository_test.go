package database_test

import (
	"github.com/ThamirisMonteiro/rate-limiter/internal/infra/database"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type RedisRepositoryTestSuite struct {
	suite.Suite
	repository database.Repository
}

func (suite *RedisRepositoryTestSuite) SetupTest() {
	suite.repository = database.NewRedisRepository("localhost:6379")
	suite.NoError(suite.repository.FlushDB())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(RedisRepositoryTestSuite))
}

func (suite *RedisRepositoryTestSuite) Test_SetKey() {
	err := suite.repository.SetKey("key", "value", 0)
	suite.NoError(err)

	value, err := suite.repository.GetKey("key")
	suite.NoError(err)
	suite.Equal("value", value)
}

func (suite *RedisRepositoryTestSuite) Test_IncrCounter() {
	err := suite.repository.IncrCounter("counter")
	suite.NoError(err)

	value, err := suite.repository.GetKey("counter")
	suite.NoError(err)
	suite.Equal("1", value)
}

func (suite *RedisRepositoryTestSuite) Test_ExpireKey() {
	err := suite.repository.SetKey("key", "value", 1*time.Second)
	suite.NoError(err)

	err = suite.repository.ExpireKey("key", 1*time.Second)
	suite.NoError(err)

	time.Sleep(time.Second)

	value, err := suite.repository.GetKey("key")
	suite.Error(err)
	suite.Empty(value)
}

func (suite *RedisRepositoryTestSuite) Test_TTLKey() {
	err := suite.repository.SetKey("key", "value", 900*time.Second)
	suite.NoError(err)

	ttl, err := suite.repository.TTLKey("key")
	suite.NoError(err)
	suite.Equal(900*time.Second, ttl)
}

func (suite *RedisRepositoryTestSuite) Test_TTLKey_Expired() {
	err := suite.repository.SetKey("key", "value", 1*time.Second)
	suite.NoError(err)

	time.Sleep(2 * time.Second)

	ttl, err := suite.repository.TTLKey("key")
	suite.NoError(err)
	suite.Equal(KeyDoesNotExist, ttl)
}

func (suite *RedisRepositoryTestSuite) Test_TTLKey_NotFound() {
	ttl, err := suite.repository.TTLKey("not_found")
	suite.NoError(err)
	suite.Equal(KeyDoesNotExist, ttl)
}

func (suite *RedisRepositoryTestSuite) Test_GetKey_NotFound() {
	_, err := suite.repository.GetKey("not_found")
	suite.Error(err)
}
