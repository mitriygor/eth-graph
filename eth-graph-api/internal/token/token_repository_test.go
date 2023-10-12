package token

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockGraphClient struct {
	mock.Mock
}

func (m *MockGraphClient) Query(ctx context.Context, query interface{}, variables map[string]interface{}) error {
	args := m.Called(ctx, query, variables)
	return args.Error(0)
}

func TestGetPoolsByToken(t *testing.T) {
	mockClient := new(MockGraphClient)

	repo := NewTokenRepository(mockClient)

	token := "test-token"
	first := 5

	mockClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(1).(*struct {
			Pools []Pool "graphql:\"pools(first: $first, where: { or: [{ token0: $token }, { token1: $token }] })\""
		})
		arg.Pools = []Pool{{ID: "pool1"}, {ID: "pool2"}}
	})

	pools, err := repo.GetPoolsByToken(context.Background(), token, first)

	assert.NoError(t, err)
	assert.Len(t, pools, 2)
	assert.Equal(t, "pool1", pools[0].ID)
	assert.Equal(t, "pool2", pools[1].ID)

	mockClient.AssertExpectations(t)
}

func TestGetVolume(t *testing.T) {
	mockClient := new(MockGraphClient)
	repo := NewTokenRepository(mockClient)
	ctx := context.TODO()
	token := "0xToken"
	from := int64(1633036800)
	to := int64(1633123200)

	mockClient.On("Query", ctx, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {

		arg := args.Get(1).(*struct {
			TokenDayDatas []TokenDayData `graphql:"tokenDayDatas(where: {token: $token, date_gte: $from, date_lte: $to})"`
		})
		arg.TokenDayDatas = []TokenDayData{
			{VolumeUSD: "100.10"},
			{VolumeUSD: "200.20"},
		}
	})

	volume, err := repo.GetVolume(ctx, token, from, to)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(volume))
	assert.Equal(t, "100.10", volume[0].VolumeUSD)
	assert.Equal(t, "200.20", volume[1].VolumeUSD)

	mockClient.AssertExpectations(t)
}
