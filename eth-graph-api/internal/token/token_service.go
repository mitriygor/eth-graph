package token

import (
	"context"
	"errors"
	"eth-graph-api/pkg/calc"
	"eth-graph-api/pkg/logger"
	"eth-graph-api/pkg/validator"
	"strconv"
)

// Service is an interface that defines contracts for interacting
// with token data, ensuring implementations provide methods for
// retrieving token pools and volume.
type Service interface {
	GetPoolsByTokenService(ctx context.Context, token string, firstStr string) ([]Pool, error)
	GetVolumeService(ctx context.Context, token string, from string, to string) (*Volume, error)
}

// tokenService is a concrete implementation of the Service interface,
// facilitating retrieval of token-related data using a Repository.
type tokenService struct {
	tokenRepo Repository
}

// NewTokenService constructs a new instance of tokenService, using
// the provided Repository `repo` to fetch data.
func NewTokenService(repo Repository) Service {
	return &tokenService{
		tokenRepo: repo,
	}
}

// GetPoolsByTokenService retrieves pools by token, ensuring
// token and firstStr inputs are validated and parsed correctly,
// executing within `ctx` context. It returns a slice of Pool or an error.
func (s *tokenService) GetPoolsByTokenService(ctx context.Context, token string, firstStr string) ([]Pool, error) {
	if !validator.IsValidToken(token) {
		return nil, errors.New("invalid token")
	}

	firstNum, err := strconv.Atoi(firstStr)
	if err != nil {
		firstNum = 5
	}

	if !validator.IsValidFirst(firstNum) {
		firstNum = 5
	}

	return s.tokenRepo.GetPoolsByToken(ctx, token, firstNum)
}

// GetVolumeService retrieves token volume data within the specified range,
// ensuring token, fromStr, and toStr inputs are validated and parsed correctly,
// executing within `ctx` context. It returns a Volume struct or an error.
func (s *tokenService) GetVolumeService(ctx context.Context, token string, fromStr string, toStr string) (*Volume, error) {

	if !validator.IsValidToken(token) {
		return nil, errors.New("invalid token")
	}

	if !validator.IsValidRange(fromStr, toStr) {
		return nil, errors.New("invalid range")
	}

	from, _ := strconv.ParseInt(fromStr, 10, 64)
	to, _ := strconv.ParseInt(toStr, 10, 64)

	tokenVolumes, err := s.tokenRepo.GetVolume(ctx, token, from, to)
	if err != nil {
		return nil, errors.New("issue to get token volume")
	}

	if len(tokenVolumes) == 0 {
		return nil, errors.New("no token volume")
	}

	var volumes []string
	for _, v := range tokenVolumes {
		volumes = append(volumes, v.VolumeUSD)
	}

	sumStr, err := calc.SumNumbers(volumes)
	if err != nil {
		logger.Error("GetVolumeService error", "error", err)
		return nil, errors.New("issue to get sum of volume")
	}

	volume := &Volume{
		Volume: sumStr,
	}

	return volume, nil
}
