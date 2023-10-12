package token

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

func (m *MockRepository) GetPoolsByToken(ctx context.Context, token string, first int) ([]Pool, error) {
	args := m.Called(ctx, token, first)
	return args.Get(0).([]Pool), args.Error(1)
}

func (m *MockRepository) GetVolume(ctx context.Context, token string, from int64, to int64) ([]TokenDayData, error) {
	args := m.Called(ctx, token, from, to)
	return args.Get(0).([]TokenDayData), args.Error(1)
}

func TestGetPoolsByTokenService(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewTokenService(mockRepo)
	ctx := context.TODO()

	cases := []struct {
		name          string
		token         string
		firstStr      string
		setupMocks    func()
		expectedPools []Pool
		expectedErr   error
	}{
		{
			name:     "get pools by token with valid input",
			token:    "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			firstStr: "5",
			setupMocks: func() {
				mockRepo.On("GetPoolsByToken", ctx, "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", 5).Return([]Pool{{ID: "0x0000d36ab86d213c14d93cd5ae78615a20596505"}, {ID: "0x0001fcbba8eb491c3ccfeddc5a5caba1a98c4c28"}}, nil).Once()
			},
			expectedPools: []Pool{{ID: "0x0000d36ab86d213c14d93cd5ae78615a20596505"}, {ID: "0x0001fcbba8eb491c3ccfeddc5a5caba1a98c4c28"}},
			expectedErr:   nil,
		},
		{
			name:          "get pools by token with invalid token",
			token:         "invalidToken",
			firstStr:      "5",
			setupMocks:    func() {},
			expectedPools: nil,
			expectedErr:   errors.New("invalid token"),
		},
		{
			name:     "get pools by token with non-numeric first",
			token:    "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			firstStr: "nonNumeric",
			setupMocks: func() {
				mockRepo.On("GetPoolsByToken", ctx, "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", 5).Return([]Pool{{ID: "0x0000d36ab86d213c14d93cd5ae78615a20596505"}, {ID: "0x0001fcbba8eb491c3ccfeddc5a5caba1a98c4c28"}}, nil).Once()
			},
			expectedPools: []Pool{{ID: "0x0000d36ab86d213c14d93cd5ae78615a20596505"}, {ID: "0x0001fcbba8eb491c3ccfeddc5a5caba1a98c4c28"}},
			expectedErr:   nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks()
			returnedPools, err := service.GetPoolsByTokenService(ctx, tc.token, tc.firstStr)
			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedPools, returnedPools)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetVolumeService(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name        string
		token       string
		fromStr     string
		toStr       string
		mockRepoFn  func(m *MockRepository)
		expectErr   bool
		expectValue string
	}{
		{
			name:    "valid case",
			token:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			fromStr: "1633036800",
			toStr:   "1633123200",
			mockRepoFn: func(m *MockRepository) {
				m.On("GetVolume", mock.Anything, "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", int64(1633036800), int64(1633123200)).
					Return([]TokenDayData{{VolumeUSD: "300.0000000000"}}, nil).Once()
			},
			expectErr:   false,
			expectValue: "300.0000000000",
		},
		{
			name:    "invalid token",
			token:   "",
			fromStr: "1633036800",
			toStr:   "1633123200",
			mockRepoFn: func(m *MockRepository) {
			},
			expectErr: true,
		},
		{
			name:    "repository error",
			token:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			fromStr: "1633036800",
			toStr:   "1633123200",
			mockRepoFn: func(m *MockRepository) {
				m.On("GetVolume", mock.Anything, "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", int64(1633036800), int64(1633123200)).
					Return([]TokenDayData{}, errors.New("mock error")).Once()
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			if tc.mockRepoFn != nil {
				tc.mockRepoFn(mockRepo)
			}

			svc := &tokenService{tokenRepo: mockRepo}

			volume, err := svc.GetVolumeService(ctx, tc.token, tc.fromStr, tc.toStr)

			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectValue, volume.Volume)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
