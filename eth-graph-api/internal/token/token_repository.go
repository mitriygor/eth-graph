package token

import (
	"context"
	"eth-graph-api/pkg/logger"
	"github.com/shurcooL/graphql"
)

// GraphClient Defines the GraphClient interface which encapsulates the ability
// to perform a GraphQL query, with its signature implying it takes
// a query `q` and a map of `variables`, executing it within a certain context `ctx`.
// Query sends a GraphQL query to the endpoint.
// It accepts a context (for timeout/aborting, etc.),
// a query structure `q` (to fill with returned data), and
// a map of variables to use in the GraphQL request.
type GraphClient interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

// Repository represents an interface defining contracts for retrieving token data.
// Implementers should provide mechanisms to retrieve data about token pools and volume.
type Repository interface {
	GetPoolsByToken(ctx context.Context, token string, first int) ([]Pool, error)
	GetVolume(ctx context.Context, token string, from int64, to int64) ([]TokenDayData, error)
}

// tokenRepository is a structure that implements the Repository interface,
// allowing retrieval of token data using a GraphQL client.
type tokenRepository struct {
	graphClient GraphClient
}

// NewTokenRepository constructs a new tokenRepository with the provided GraphQL client.
func NewTokenRepository(graphClient GraphClient) Repository {
	return &tokenRepository{
		graphClient: graphClient,
	}
}

// GetPoolsByToken performs a GraphQL query to retrieve pools associated
// with the given `token`, and limited by the `first` parameter, executing within `ctx` context.
// It returns slices of Pool or an error if the query operation fails.
func (tr *tokenRepository) GetPoolsByToken(ctx context.Context, token string, first int) ([]Pool, error) {

	var query struct {
		Pools []Pool `graphql:"pools(first: $first, where: { or: [{ token0: $token }, { token1: $token }] })"`
	}

	vars := map[string]interface{}{
		"token": graphql.String(token),
		"first": graphql.Int(first),
	}

	err := tr.graphClient.Query(ctx, &query, vars)
	if err != nil {
		logger.Error("GetPoolsByToken error", "error", err)
		return nil, err
	}

	return query.Pools, nil
}

// GetVolume performs a GraphQL query to retrieve volume data associated
// with the given `token` between `from` and `to` timestamps, executing within `ctx` context.
// It returns slices of TokenDayData or an error if the query operation fails.
func (tr *tokenRepository) GetVolume(ctx context.Context, token string, from int64, to int64) ([]TokenDayData, error) {

	var query struct {
		TokenDayDatas []TokenDayData `graphql:"tokenDayDatas(where: {token: $token, date_gte: $from, date_lte: $to})"`
	}

	vars := map[string]interface{}{
		"token": graphql.String(token),
		"from":  graphql.Int(from),
		"to":    graphql.Int(to),
	}

	err := tr.graphClient.Query(ctx, &query, vars)
	if err != nil {
		logger.Error("GetVolume error", "error", err)
		return nil, err
	}

	return query.TokenDayDatas, nil
}
