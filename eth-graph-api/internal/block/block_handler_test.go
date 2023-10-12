package block

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockBlockService struct {
	mock.Mock
}

func (m *MockBlockService) GetSwapsByBlockService(ctx context.Context, blockStr string, firstStr string) ([]Swap, error) {
	args := m.Called(ctx, blockStr, firstStr)
	return args.Get(0).([]Swap), args.Error(1)
}

func (m *MockBlockService) GetSwappedTokensByBlockService(ctx context.Context, blockStr string, firstStr string) ([]Token, error) {
	args := m.Called(ctx, blockStr, firstStr)
	return args.Get(0).([]Token), args.Error(1)
}

func TestGetSwapsByBlockHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockSvcOutput  []Swap
		mockSvcErr     error
		expectedStatus int
	}{
		{
			name:           "valid request",
			url:            "/v1/blocks/18319881/swaps?first=5",
			mockSvcOutput:  []Swap{{ID: "1"}},
			mockSvcErr:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid block",
			url:            "/v1/blocks/invalid/swaps?first=5",
			mockSvcOutput:  nil,
			mockSvcErr:     errors.New("invalid block"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockSvc := new(MockBlockService)
			mockSvc.On("GetSwapsByBlockService", mock.Anything, mock.Anything, mock.Anything).Return(test.mockSvcOutput, test.mockSvcErr).Once()

			h := &Handler{
				BlockService: mockSvc,
			}

			req, err := http.NewRequest(http.MethodGet, test.url, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Route("/v1/blocks/{block}/swaps", func(r chi.Router) {
				r.Get("/", h.GetSwapsByBlockHandler)
			})
			r.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestGetSwappedTokensByBlockHandler(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockSvcOutput  []Token
		mockSvcErr     error
		expectedStatus int
	}{
		{
			name:           "valid request",
			url:            "/v1/blocks/18319881/swaps/tokens?first=5",
			mockSvcOutput:  []Token{{Symbol: "ETH"}},
			mockSvcErr:     nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid block",
			url:            "/v1/blocks/invalid/swaps/tokens?first=5",
			mockSvcOutput:  nil,
			mockSvcErr:     errors.New("invalid block"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockSvc := new(MockBlockService)
			mockSvc.On("GetSwappedTokensByBlockService", mock.Anything, mock.Anything, mock.Anything).Return(test.mockSvcOutput, test.mockSvcErr).Once()

			h := &Handler{
				BlockService: mockSvc,
			}

			req, err := http.NewRequest(http.MethodGet, test.url, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Route("/v1/blocks/{block}/swaps/tokens", func(r chi.Router) {
				r.Get("/", h.GetSwappedTokensByBlockHandler)
			})
			r.ServeHTTP(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)

			mockSvc.AssertExpectations(t)
		})
	}
}
