package json_helper

import (
	"encoding/json"
	"net/http"
)

// jsonResponse struct is used for marshaling the JSON response payload.
// It includes:
// - Error: a boolean that is true when an error occurs.
// - Message: a string that may hold error details or other messages.
// - Data: a field that can hold any type of data and will be omitted if empty.
type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// WriteJSON writes a JSON response with provided HTTP status, data,
// and optional headers to an http.ResponseWriter.
//
// Parameters:
// - `w`: the http.ResponseWriter where the response will be written.
// - `status`: the HTTP status code to be returned in the response.
// - `data`: the payload to be serialized into JSON format and written in the response body.
// - `headers` (optional): additional HTTP headers to be included in the response.
//
// Returns:
// - An error, if serialization to JSON fails or if the response cannot be written.
//
// Example usage:
//
//	err := WriteJSON(w, http.StatusOK, someData)
//	if err != nil {
//	    log.Printf("Could not write JSON response: %v", err)
//	}
func WriteJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// ErrorJSON writes an error message as JSON format with an optional HTTP status.
// If no status is provided, it defaults to http.StatusBadRequest (400).
//
// Parameters:
// - `w`: the http.ResponseWriter where the response will be written.
// - `err`: the error to be written in the response's message field.
// - `status` (optional): the HTTP status code to be returned in the response.
//
// Returns:
// - An error, if writing the JSON response fails.
//
// Example usage:
//
//	if someError != nil {
//	    err := ErrorJSON(w, someError, http.StatusInternalServerError)
//	    if err != nil {
//	        log.Printf("Could not write error JSON response: %v", err)
//	    }
//	}
func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return WriteJSON(w, statusCode, payload)
}
