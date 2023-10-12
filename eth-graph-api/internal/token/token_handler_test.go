package token

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetVolumeService(ctx context.Context, token, fromStr, toStr string) (*Volume, error) {
	args := m.Called(ctx, token, fromStr, toStr)
	return args.Get(0).(*Volume), args.Error(1)
}

func (m *MockService) GetPoolsByTokenService(ctx context.Context, token string, firstStr string) ([]Pool, error) {
	args := m.Called(ctx, token, firstStr)
	return args.Get(0).([]Pool), args.Error(1)
}

func TestGetPoolsByTokenHandler(t *testing.T) {
	mockService := new(MockService)

	handler := &Handler{
		TokenService: mockService,
	}

	tests := []struct {
		name           string
		token          string
		mockServiceFn  func()
		expectedStatus int
	}{
		{
			name:  "success",
			token: "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			mockServiceFn: func() {
				mockService.On("GetPoolsByTokenService", mock.Anything, "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", "5").
					Return([]Pool{{ID: "0x0000d36ab86d213c14d93cd5ae78615a20596505"}}, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/v1/tokens/"+tt.token+"/pools", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			if tt.mockServiceFn != nil {
				tt.mockServiceFn()
			}

			r := chi.NewRouter()
			r.Get("/v1/tokens/{token}/pools", handler.GetPoolsByTokenHandler)
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetVolumeHandler(t *testing.T) {
	mockService := new(MockService)

	handler := &Handler{
		TokenService: mockService,
	}

	tests := []struct {
		name           string
		token          string
		mockServiceFn  func()
		expectedStatus int
	}{
		{
			name: "success",
			mockServiceFn: func() {
				mockService.On("GetVolumeService", mock.Anything, "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", mock.Anything, mock.Anything).
					Return(&Volume{Volume: "8039962301.3240150269"}, nil)
			},
			token:          "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/v1/tokens/"+tt.token+"/volume", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()

			if tt.mockServiceFn != nil {
				tt.mockServiceFn()
			}

			r := chi.NewRouter()
			r.Get("/v1/tokens/{token}/volume", handler.GetVolumeHandler)
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			mockService.AssertExpectations(t)
		})
	}
}
