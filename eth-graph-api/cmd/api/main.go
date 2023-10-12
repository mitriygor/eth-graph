package main

import (
	"eth-graph-api/pkg/logger"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/shurcooL/graphql"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

func main() {
	// Initialize the logger
	err := logger.Initialize()
	if err != nil {
		panic(err)
	}
	defer func(Log *zap.SugaredLogger) {
		err := Log.Sync()
		if err != nil {
			log.Printf("Error syncing logger: %v", err)
		}
	}(logger.Log)

	// Load environment variables from a .env file.
	// godotenv.Load() returns an error if it cannot find or load the .env file.
	err = godotenv.Load()
	if err != nil {
		logger.Error("error occurred", "error", err)
	}

	// Get the port number and Graph API URL from environment variables.
	port := os.Getenv("PORT")
	graphApi := os.Getenv("GRAPH_API")
	apiVersion := fmt.Sprintf("/%s", os.Getenv("API_VERSION"))

	// Create a new GraphQL client using the API URL.
	graphqlClient := graphql.NewClient(graphApi, nil)

	// Initialize the API handlers using the API version and GraphQL client
	// The initHandlers function is presumed to set up the routing and handlers for the API,
	// and set up the HTTP server using the specified port and handler.
	router := initHandlers(apiVersion, graphqlClient)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}

	// Start the HTTP server and listen for requests
	// ListenAndServe returns an error if the server closes abnormally
	err = srv.ListenAndServe()
	if err != nil {
		logger.Error("error occurred", "error", err)
	}
}
