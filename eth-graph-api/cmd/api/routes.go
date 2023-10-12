package main

import (
	"eth-graph-api/internal/block"
	"eth-graph-api/internal/token"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"time"
)

// Routes initializes and returns an http.Handler that handles routing for the API.
// It uses the given apiVersion to prefix the API routes and uses the provided
// tokenHandler and blockHandler to handle requests to token- and block-related routes, respectively
func Routes(apiVersion string, tokenHandler token.Handler, blockHandler block.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Every(1*time.Second), 1)
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(requestLogger)
	mux.Use(rateLimit(limiter))

	mux.Route(apiVersion, func(mux chi.Router) {
		mux.Route("/tokens/{token}", func(mux chi.Router) {
			mux.Get("/pools", tokenHandler.GetPoolsByTokenHandler)
			mux.Get("/volume", tokenHandler.GetVolumeHandler)
		})
		mux.Route("/blocks/{block}", func(mux chi.Router) {
			mux.Get("/swaps", blockHandler.GetSwapsByBlockHandler)
			mux.Get("/swaps/tokens", blockHandler.GetSwappedTokensByBlockHandler)
		})
	})

	return mux
}

// rateLimit returns a middleware function that applies rate limiting to incoming HTTP requests
// based on the provided rate.Limiter. If the rate limit is exceeded, the middleware returns
// a 429 Too Many Requests error
func rateLimit(limiter *rate.Limiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// requestLogger is a middleware function that logs the details of incoming HTTP requests and
// the corresponding responses, including method, URI, status code, and duration. It uses
// zap to log the details
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, _ := zap.NewProduction()
		defer func(logger *zap.Logger) {
			err := logger.Sync()
			if err != nil {
				log.Printf("Error syncing logger: %v", err)
			}
		}(logger)

		start := time.Now()
		logger.Info("Request",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
		)

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		logger.Info("Response",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Int("status", ww.Status()),
			zap.Duration("duration", time.Since(start)),
		)
	})
}
