package ratelimiter

import (
	"github.com/ThamirisMonteiro/rate-limiter/internal/infra/database"
	"strconv"
	"time"
)

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

	err := service.repo.SetKey("req_limit", strconv.Itoa(reqLimit), 0)
	if err != nil {
		return nil, err
	}

	err = service.repo.SetKey("block_time_ip", blockTimeIP.String(), 0)
	if err != nil {
		return nil, err
	}

	err = service.repo.SetKey("block_time_token", blockTimeToken.String(), 0)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (s Service) AllowRequest(identifier, reqType string) (bool, error) {
	exists, err := s.AlreadyExists(identifier)
	if err != nil {
		return false, err
	}

	if !exists {
		var blockTime int
		if s.blockTimeToken > s.blockTimeIP {
			blockTime = int(s.blockTimeToken)
		} else {
			if reqType == "ip" {
				blockTime = int(s.blockTimeIP)
			} else {
				blockTime = int(s.blockTimeToken)
			}
		}

		err := s.repo.SetKey("requests:"+identifier, "1", time.Duration(blockTime))
		if err != nil {
			return false, err
		}

	} else {
		count, err := s.GetRequestCount(identifier)
		if err != nil {
			return false, err
		}

		ttl, err := s.GetTTL(identifier)
		if err != nil {
			return false, err
		}

		if count >= s.reqLimit || ttl <= 0 {
			return false, nil
		} else {
			err := s.repo.IncrCounter("requests:" + identifier)
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
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

func (s Service) GetTTL(identifier string) (int, error) {
	ttl, err := s.repo.TTLKey("requests:" + identifier)
	if err != nil {
		return 0, err
	}

	return int(ttl), nil
}

func (s Service) AlreadyExists(identifier string) (bool, error) {
	_, err := s.repo.GetKey("requests:" + identifier)
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}

		return false, err
	}
	return true, nil
}
