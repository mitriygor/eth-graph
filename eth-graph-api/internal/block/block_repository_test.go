package block

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockGraphClient struct {
	mock.Mock
}

func (m *MockGraphClient) Query(ctx context.Context, query interface{}, vars map[string]interface{}) error {
	args := m.Called(ctx, query, vars)
	return args.Error(0)
}

func TestGetSwapsByBlock(t *testing.T) {
	tests := []struct {
		name           string
		block          int
		first          int
		mockClientFunc func(m *MockGraphClient)
		expectedResult []Swap
		expectedError  string
	}{
		{
			name:  "success",
			block: 18319881,
			first: 5,
			mockClientFunc: func(m *MockGraphClient) {
				m.On("Query", mock.Anything, mock.Anything, mock.Anything).
					Return(nil).
					Run(func(args mock.Arguments) {
						arg := args.Get(1).(*struct {
							Swaps []Swap `graphql:"swaps(block: {number: $block}, first: $first)"`
						})
						arg.Swaps = []Swap{{ID: "1"}, {ID: "2"}}
					})
			},
			expectedResult: []Swap{{ID: "1"}, {ID: "2"}},
		},
		{
			name:  "error from client",
			block: 123456,
			first: 10,
			mockClientFunc: func(m *MockGraphClient) {
				m.On("Query", mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("client error"))
			},
			expectedError: "client error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := new(MockGraphClient)
			test.mockClientFunc(mockClient)

			repo := NewBlockRepository(mockClient)
			result, err := repo.GetSwapsByBlock(context.Background(), test.block, test.first)

			if test.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, test.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedResult, result)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
