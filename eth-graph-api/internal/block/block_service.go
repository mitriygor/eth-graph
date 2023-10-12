package block

import (
	"context"
	"errors"
	"eth-graph-api/pkg/validator"
	"strconv"
)

// Service is an interface that declares methods for retrieving swap and token data.
// Implementations of this interface will provide concrete functionality
// for fetching and possibly processing swap and token data.
type Service interface {
	GetSwapsByBlockService(ctx context.Context, blockStr string, firstStr string) ([]Swap, error)
	GetSwappedTokensByBlockService(ctx context.Context, blockStr string, firstStr string) ([]Token, error)
}

// blockService is a struct that implements the Service interface.
// It uses a Repository to fetch swap data and implements additional logic
// to process and return the requested data.
type blockService struct {
	blockRepo Repository
}

// NewBlockService is a constructor function that creates and returns a new instance
// of the blockService, initializing it with a provided Repository.
func NewBlockService(repo Repository) Service {
	return &blockService{
		blockRepo: repo,
	}
}

// GetSwapsByBlockService retrieves swap data for a specific blockchain block.
// - It validates and converts input parameters (block and first) from strings to integers,
// - Fetches swap data using the blockRepo,
// - And returns the fetched data.
func (s *blockService) GetSwapsByBlockService(ctx context.Context, blockStr string, firstStr string) ([]Swap, error) {

	if !validator.IsValidBlock(blockStr) {
		return nil, errors.New("invalid block")
	}

	firstNum, err := strconv.Atoi(firstStr)
	if err != nil {
		firstNum = 5
	}

	if !validator.IsValidFirst(firstNum) {
		firstNum = 5
	}

	blockNum, err := strconv.Atoi(blockStr)
	if err != nil {
		return nil, errors.New("invalid block")
	}

	return s.blockRepo.GetSwapsByBlock(ctx, blockNum, firstNum)
}

// GetSwappedTokensByBlockService retrieves token data for swaps in a specific blockchain block.
// - It validates and converts input parameters (block and first) from strings to integers,
// - Retrieves swap data using the blockRepo,
// - Processes swaps to extract and collate related token data,
// - And returns the token data.
func (s *blockService) GetSwappedTokensByBlockService(ctx context.Context, blockStr string, firstStr string) ([]Token, error) {

	if !validator.IsValidBlock(blockStr) {
		return nil, errors.New("invalid block")
	}

	firstNum, err := strconv.Atoi(firstStr)
	if err != nil {
		firstNum = 5
	}

	if !validator.IsValidFirst(firstNum) {
		firstNum = 5
	}

	blockNum, err := strconv.Atoi(blockStr)
	if err != nil {
		return nil, errors.New("invalid block")
	}

	swaps, err := s.blockRepo.GetSwapsByBlock(ctx, blockNum, firstNum)

	if err != nil {
		return nil, errors.New("issue to get swaps")
	}

	var tokensMap = make(map[string]Token)
	var tokens []Token

	for _, s := range swaps {
		if _, ok := tokensMap[s.Pool.Token0.Symbol]; !ok {
			tokensMap[s.Pool.Token0.Symbol] = s.Pool.Token0
		}

		if _, ok := tokensMap[s.Pool.Token1.Symbol]; !ok {
			tokensMap[s.Pool.Token1.Symbol] = s.Pool.Token1
		}
	}

	for _, t := range tokensMap {
		tokens = append(tokens, t)
	}

	return tokens, nil
}
