package ratelimiter

import (
	"github.com/ThamirisMonteiro/rate-limiter/internal/infra/database"
	"strconv"
	"time"
)

type Service struct {
	Repo           database.Repository
	ReqLimit       int
	BlockTimeIP    time.Duration
	BlockTimeToken time.Duration
}

func NewService(repo database.Repository, reqLimit int, blockTimeIP, blockTimeToken time.Duration) (*Service, error) {
	service := &Service{
		Repo:           repo,
		ReqLimit:       reqLimit,
		BlockTimeIP:    blockTimeIP,
		BlockTimeToken: blockTimeToken,
	}

	err := service.Repo.SetKey("req_limit", strconv.Itoa(reqLimit), 0)
	if err != nil {
		return nil, err
	}

	err = service.Repo.SetKey("block_time_ip", blockTimeIP.String(), 0)
	if err != nil {
		return nil, err
	}

	err = service.Repo.SetKey("block_time_token", blockTimeToken.String(), 0)
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
		if s.BlockTimeToken > s.BlockTimeIP {
			blockTime = int(s.BlockTimeToken)
		} else {
			if reqType == "ip" {
				blockTime = int(s.BlockTimeIP)
			} else {
				blockTime = int(s.BlockTimeToken)
			}
		}

		err := s.Repo.SetKey("requests:"+identifier, "1", time.Duration(blockTime))
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

		if count >= s.ReqLimit || ttl <= 0 {
			return false, nil
		} else {
			err := s.Repo.IncrCounter("requests:" + identifier)
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func (s Service) GetRequestCount(identifier string) (int, error) {
	counterKey := "requests:" + identifier

	countStr, err := s.Repo.GetKey(counterKey)
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
	ttl, err := s.Repo.TTLKey("requests:" + identifier)
	if err != nil {
		return 0, err
	}

	return int(ttl), nil
}

func (s Service) AlreadyExists(identifier string) (bool, error) {
	_, err := s.Repo.GetKey("requests:" + identifier)
	if err != nil {
		if err.Error() == "redis: nil" {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
