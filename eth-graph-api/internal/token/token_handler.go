package token

import (
	"context"
	"errors"
	jh "eth-graph-api/pkg/json_helper"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

// Handler is a struct that contains a Service which provides
// methods for token data retrieval.
type Handler struct {
	TokenService Service
}

// GetPoolsByTokenHandler is an HTTP handler function that retrieves pool data for a specific token.
// The token parameter is extracted from the URL,
// and optionally a 'first' query parameter can specify the maximum number of results to return.
// It responds with JSON-encoded pool data or appropriate error responses.
func (h *Handler) GetPoolsByTokenHandler(w http.ResponseWriter, r *http.Request) {

	token := chi.URLParam(r, "token")
	if token == "" {
		err := jh.ErrorJSON(w, errors.New("token cannot be empty"), http.StatusBadRequest)
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

	pools, err := h.TokenService.GetPoolsByTokenService(ctx, token, firstStr)
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

	err = jh.WriteJSON(w, http.StatusOK, pools)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetVolumeHandler is an HTTP handler function that retrieves token volume data.
// It extracts 'token', 'from', and 'to' parameters from the request,
// then utilizes TokenService to retrieve and respond with the volume data,
// or handle errors appropriately.
func (h *Handler) GetVolumeHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		err := jh.ErrorJSON(w, errors.New("token cannot be empty"), http.StatusBadRequest)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	queryParams := r.URL.Query()
	now := time.Now()
	nowTimestamp := now.Unix()

	toStr := queryParams.Get("to")
	if toStr == "" {
		toStr = strconv.FormatInt(nowTimestamp, 10)
	}

	fromStr := queryParams.Get("from")
	if fromStr == "" {
		yesterday := now.Add(-24 * time.Hour)
		yesterdayTimestamp := yesterday.Unix()
		fromStr = strconv.FormatInt(yesterdayTimestamp, 10)
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	volume, err := h.TokenService.GetVolumeService(ctx, token, fromStr, toStr)
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

	err = jh.WriteJSON(w, http.StatusOK, volume)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
