package main

import (
	"context"
	"eth-graph-api/internal/block"
	"eth-graph-api/internal/token"
	"github.com/shurcooL/graphql"
	"net/http"
)

// RealGraphClient is a structure that wraps the graphql.Client
// to make it conform to our application-specific interface, allowing
// it to be used to perform GraphQL queries
type RealGraphClient struct {
	Client *graphql.Client
}

// Query performs a GraphQL query using the embedded graphql.Client,
// with the provided context, query structure, and variables.
// It returns an error if the query operation fails
func (r *RealGraphClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return r.Client.Query(ctx, q, variables)
}

// initHandlers initializes and returns the HTTP handlers for the API
// given a particular API version and GraphQL client. It also sets
// up the repositories, services, and handlers for the token and block
// resources
func initHandlers(apiVersion string, graphClient *graphql.Client) http.Handler {

	tokenRepo := token.NewTokenRepository(&RealGraphClient{Client: graphClient})
	tokenService := token.NewTokenService(tokenRepo)
	tokenHandler := &token.Handler{TokenService: tokenService}

	blockRepo := block.NewBlockRepository(&RealGraphClient{Client: graphClient})
	blockService := block.NewBlockService(blockRepo)
	blockHandler := &block.Handler{BlockService: blockService}

	return Routes(apiVersion, *tokenHandler, *blockHandler)
}
