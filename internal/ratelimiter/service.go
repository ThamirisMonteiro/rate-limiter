package ratelimiter

import (
	"github.com/ThamirisMonteiro/rate-limiter/internal/infra/database"
	"log"
	"os"
	"strconv"
	"time"
)

var SettingsKey = "rate_limit:settings"

type Service struct {
	repository database.Repository
	reqLimit   int
}

func NewService(repository database.Repository) *Service {
	reqLimitStr := os.Getenv("REQ_LIMIT")

	reqLimit, err := strconv.Atoi(reqLimitStr)
	if err != nil {
		log.Printf("Invalid REQ_LIMIT value: %s. Defaulting to 0.", reqLimitStr)
		reqLimit = 0
	}

	return &Service{
		repository: repository,
		reqLimit:   reqLimit,
	}
}

func (s Service) AllowRequest(identifier string) (bool, error) {
	currentCount, err := s.GetRequestCount(identifier)
	if err != nil {
		return false, err
	}

	if currentCount < s.reqLimit {
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
