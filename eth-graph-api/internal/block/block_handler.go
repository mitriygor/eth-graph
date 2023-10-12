package block

import (
	"context"
	"errors"
	jh "eth-graph-api/pkg/json_helper"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type Handler struct {
	BlockService Service
}

// GetSwapsByBlockHandler is an HTTP handler function that retrieves
// and sends swaps data for a specific block in response to HTTP requests.
// - It extracts the "block" URL parameter and validates it.
// - Parses and validates query parameters.
// - Sets a 5-second timeout for the request context.
// - Fetches swaps data for the block using BlockService.
// - Handles potential errors and timeouts.
// - Writes the fetched data as JSON to the HTTP response.
func (h *Handler) GetSwapsByBlockHandler(w http.ResponseWriter, r *http.Request) {

	block := chi.URLParam(r, "block")
	if block == "" {
		err := jh.ErrorJSON(w, errors.New("block cannot be empty"), http.StatusBadRequest)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	queryParams := r.URL.Query()
	firstStr := queryParams.Get("first")
	if firstStr == "" {
		firstStr = "5"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	swaps, err := h.BlockService.GetSwapsByBlockService(ctx, block, firstStr)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = jh.ErrorJSON(w, errors.New("request timeout"), http.StatusRequestTimeout)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		} else {
			err = jh.ErrorJSON(w, err, http.StatusBadRequest)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}
		return
	}

	err = jh.WriteJSON(w, http.StatusOK, swaps)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetSwappedTokensByBlockHandler is an HTTP handler function that retrieves
// and sends swapped tokens data for a specific block in response to HTTP requests.
// - It extracts the "block" URL parameter and validates it.
// - Parses and validates query parameters.
// - Sets a 5-second timeout for the request context.
// - Fetches swapped tokens data for the block using BlockService.
// - Handles potential errors and timeouts.
// - Writes the fetched data as JSON to the HTTP response.
func (h *Handler) GetSwappedTokensByBlockHandler(w http.ResponseWriter, r *http.Request) {

	block := chi.URLParam(r, "block")
	if block == "" {
		err := jh.ErrorJSON(w, errors.New("block cannot be empty"), http.StatusBadRequest)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	queryParams := r.URL.Query()
	firstStr := queryParams.Get("first")
	if firstStr == "" {
		firstStr = "5"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	tokens, err := h.BlockService.GetSwappedTokensByBlockService(ctx, block, firstStr)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = jh.ErrorJSON(w, errors.New("request timeout"), http.StatusRequestTimeout)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		} else {
			err = jh.ErrorJSON(w, err, http.StatusBadRequest)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}
		return
	}

	err = jh.WriteJSON(w, http.StatusOK, tokens)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
