package ratelimiter

import (
	"github.com/ThamirisMonteiro/rate-limiter/internal/infra/database"
	"strconv"
	"time"
)

var SettingsKey = "rate_limit:settings"

type Service struct {
	repo           database.Repository
	reqLimit       int
	blockTimeIP    time.Duration
	blockTimeToken time.Duration
}

func NewService(repo database.Repository, reqLimit int, blockTimeIP, blockTimeToken time.Duration) (*Service, error) {
	service := &Service{
		repo:           repo,
		reqLimit:       reqLimit,
		blockTimeIP:    blockTimeIP,
		blockTimeToken: blockTimeToken,
	}

	err := service.SetLimit("request_limit", reqLimit, blockTimeToken)
	if err != nil {
		return nil, err
	}

	err = service.SetLimit("block_time_ip", int(blockTimeIP.Seconds()), blockTimeIP)
	if err != nil {
		return nil, err
	}

	err = service.SetLimit("block_time_token", int(blockTimeToken.Seconds()), blockTimeToken)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s Service) AllowRequest(identifier string) (bool, error) {
	currentCount, err := s.GetRequestCount(identifier)
	if err != nil {
		return false, err
	}

	if currentCount < s.reqLimit {
		if currentCount == 0 {
			err := s.repo.IncrCounter("requests:" + identifier)
			if err != nil {
				return false, err
			}

			err = s.repo.ExpireKey("requests:"+identifier, time.Hour)
			if err != nil {
				return false, err
			}
			return true, nil
		}

		err = s.repo.IncrCounter("requests:" + identifier)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (s Service) GetRequestCount(identifier string) (int, error) {
	counterKey := "requests:" + identifier

	countStr, err := s.repo.GetKey(counterKey)
	if err != nil {
		if err.Error() == "redis: nil" {
			return 0, nil
		}

		return 0, err
	}

	if countStr == "" {
		return 0, nil
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s Service) SetLimit(key string, limit int, ttl time.Duration) error {
	err := s.repo.SetKey(key, strconv.Itoa(limit), ttl)
	if err != nil {
		return err
	}

	return nil
}
