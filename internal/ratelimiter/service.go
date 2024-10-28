package ratelimiter

import (
	"github.com/ThamirisMonteiro/rate-limiter/internal/infra/database"
	"os"
	"strconv"
	"time"
)

var SettingsKey = "rate_limit:settings"

type Service struct {
	repository database.Repository
}

func NewService(repository database.Repository) *Service {
	return &Service{repository: repository}
}

func (s Service) AllowRequest(identifier string) (bool, error) {
	reqLimitStr := os.Getenv("REQ_LIMIT")

	reqLimit, err := strconv.Atoi(reqLimitStr)
	if err != nil {
		return false, err
	}

	currentCount, err := s.GetRequestCount(identifier)
	if err != nil {
		return false, err
	}

	if currentCount < reqLimit {
		err := s.repository.IncrCounter("requests:" + identifier)
		if err != nil {
			return false, err
		}

		if currentCount == 0 {
			err := s.repository.ExpireKey("requests:"+identifier, time.Hour)
			if err != nil {
				return false, err
			}
		}

		return true, nil
	}

	return false, nil
}

func (s Service) GetRequestCount(identifier string) (int, error) {
	counterKey := "requests:" + identifier

	countStr, err := s.repository.GetKey(counterKey)
	if err != nil {
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

func (s Service) SetLimit(limit int, ttl time.Duration) error {
	err := s.repository.SetKey(SettingsKey, strconv.Itoa(limit), ttl)
	if err != nil {
		return err
	}

	err = s.repository.ExpireKey(SettingsKey, ttl)
	if err != nil {
		return err
	}

	return nil
}
