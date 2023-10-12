package block

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetSwapsByBlock(ctx context.Context, block int, first int) ([]Swap, error) {
	args := m.Called(ctx, block, first)
	return args.Get(0).([]Swap), args.Error(1)
}

func TestGetSwapsByBlockService(t *testing.T) {
	tests := []struct {
		name           string
		blockStr       string
		firstStr       string
		mockRepoOutput []Swap
		mockRepoErr    error
		expectedOutput []Swap
		expectingError bool
	}{
		{
			name:           "valid input",
			blockStr:       "1234",
			firstStr:       "10",
			mockRepoOutput: []Swap{{ID: "1"}},
			mockRepoErr:    nil,
			expectedOutput: []Swap{{ID: "1"}},
			expectingError: false,
		},
		{
			name:           "invalid block",
			blockStr:       "invalid",
			firstStr:       "10",
			mockRepoOutput: nil,
			mockRepoErr:    nil,
			expectedOutput: nil,
			expectingError: true,
		},
		{
			name:           "repo returns error",
			blockStr:       "1234",
			firstStr:       "10",
			mockRepoOutput: nil,
			mockRepoErr:    errors.New("some error"),
			expectedOutput: nil,
			expectingError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(MockRepository)

			if !test.expectingError || (test.expectingError && test.mockRepoErr != nil) {
				mockRepo.On("GetSwapsByBlock", mock.Anything, mock.Anything, mock.Anything).Return(test.mockRepoOutput, test.mockRepoErr).Once()
			}

			svc := NewBlockService(mockRepo)
			output, err := svc.GetSwapsByBlockService(context.Background(), test.blockStr, test.firstStr)

			if test.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedOutput, output)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetSwappedTokensByBlockService(t *testing.T) {
	tests := []struct {
		name           string
		blockStr       string
		firstStr       string
		mockRepoOutput []Swap
		mockRepoErr    error
		expectedOutput []Token
		expectingError bool
	}{
		{
			name:           "valid input",
			blockStr:       "1234",
			firstStr:       "10",
			mockRepoOutput: []Swap{{ID: "1", Pool: Pool{Token0: Token{Symbol: "AAA"}, Token1: Token{Symbol: "BBB"}}}},
			mockRepoErr:    nil,
			expectedOutput: []Token{{Symbol: "AAA"}, {Symbol: "BBB"}},
			expectingError: false,
		},
		{
			name:           "invalid block",
			blockStr:       "invalid",
			firstStr:       "10",
			mockRepoOutput: nil,
			mockRepoErr:    nil,
			expectedOutput: nil,
			expectingError: true,
		},
		{
			name:           "repo returns error",
			blockStr:       "1234",
			firstStr:       "10",
			mockRepoOutput: nil,
			mockRepoErr:    errors.New("some error"),
			expectedOutput: nil,
			expectingError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(MockRepository)

			if !test.expectingError || (test.expectingError && test.mockRepoErr != nil) {
				mockRepo.On("GetSwapsByBlock", mock.Anything, mock.Anything, mock.Anything).Return(test.mockRepoOutput, test.mockRepoErr).Once()
			}

			svc := NewBlockService(mockRepo)
			output, err := svc.GetSwappedTokensByBlockService(context.Background(), test.blockStr, test.firstStr)

			if test.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, test.expectedOutput, output)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
