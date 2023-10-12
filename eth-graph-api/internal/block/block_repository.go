package block

import (
	"context"
	"eth-graph-api/pkg/logger"
	"github.com/shurcooL/graphql"
)

// GraphClient is an interface that declares a method for making
// GraphQL queries, which can be implemented by various clients
// that interact with a GraphQL API.
// Query sends a GraphQL query to the API and populates the response data
// into the passed query structure. "variables" parameter is used to provide
// GraphQL variables in the query. The method returns an error if the query
// execution fails.
type GraphClient interface {
	Query(ctx context.Context, q interface{}, variables map[string]interface{}) error
}

// Repository is an interface that declares methods for fetching Swap data,
// providing a way to access Swap data without exposing details of the data retrieval.
// GetSwapsByBlock retrieves swap data related to a specific blockchain block.
// It takes a block number and a maximum number of results ("first") as parameters,
// and returns a slice of Swaps and an error if the data retrieval fails
type Repository interface {
	GetSwapsByBlock(ctx context.Context, block int, first int) ([]Swap, error)
}

// blockRepository is a struct that implements the Repository interface,
// it uses a GraphClient to fetch swap data related to blockchain blocks
// graphClient is used to execute GraphQL queries against an API.
type blockRepository struct {
	graphClient GraphClient
}

// NewBlockRepository is a constructor function that returns a new instance of
// a struct implementing the Repository interface, initializing it with
// a provided GraphClient.
func NewBlockRepository(graphClient GraphClient) Repository {
	return &blockRepository{
		graphClient: graphClient,
	}
}

// GetSwapsByBlock is a method on blockRepository that fetches and returns
// swap data for a specific blockchain block from the GraphQL API.
// - It constructs a GraphQL query and variables based on provided parameters,
// - Sends the query to the GraphQL API using the graphClient,
// - Handles potential query execution errors,
// - And returns the fetched swap data.
func (tr *blockRepository) GetSwapsByBlock(ctx context.Context, block int, first int) ([]Swap, error) {

	var query struct {
		Swaps []Swap `graphql:"swaps(block: {number: $block}, first: $first)"`
	}

	vars := map[string]interface{}{
		"block": graphql.Int(block),
		"first": graphql.Int(first),
	}

	err := tr.graphClient.Query(ctx, &query, vars)
	if err != nil {
		logger.Error("error querying swaps by block", err)
		return nil, err
	}

	return query.Swaps, nil
}
