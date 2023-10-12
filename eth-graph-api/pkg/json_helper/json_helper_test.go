package json_helper

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	data := jsonResponse{
		Error:   false,
		Message: "Success",
		Data:    "Test data",
	}

	headers := http.Header{}
	headers.Add("Custom-Header", "CustomValue")

	w := httptest.NewRecorder()
	err := WriteJSON(w, http.StatusOK, data, headers)

	if err != nil {
		t.Errorf("WriteJSON: Unexpected error: %v", err)
	}

	result := w.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("WriteJSON: Expected status OK, got: %v", result.Status)
	}

	if contentType := result.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("WriteJSON: Expected application/json, got: %v", contentType)
	}

	for key, value := range headers {
		if !reflect.DeepEqual(result.Header[key], value) {
			t.Errorf("WriteJSON: Expected %v for header key %v, got: %v", value, key, result.Header[key])
		}
	}

	var response jsonResponse
	err = json.NewDecoder(result.Body).Decode(&response)
	if err != nil {
		t.Errorf("WriteJSON: Error decoding JSON: %v", err)
	}

	if !reflect.DeepEqual(response, data) {
		t.Errorf("WriteJSON: Expected %+v, got: %+v", data, response)
	}
}

func TestErrorJSON(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
		statusCode   int
		expectedMsg  string
	}{
		{
			name:         "With default status code",
			err:          errors.New("error occurred"),
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "error occurred",
		},
		{
			name:         "With custom status code",
			err:          errors.New("not found"),
			expectedCode: http.StatusNotFound,
			statusCode:   http.StatusNotFound,
			expectedMsg:  "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := http.NewRequest("GET", "/does/not/exist", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			if tt.statusCode > 0 {
				err := ErrorJSON(rr, tt.err, tt.statusCode)
				if err != nil {
					t.Errorf("Error happened")
				}
			} else {
				err := ErrorJSON(rr, tt.err)
				if err != nil {
					t.Errorf("Error happened")
				}
			}

			if rr.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			var resp jsonResponse
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			if err != nil {
				t.Fatal(err)
			}

			if resp.Message != tt.expectedMsg {
				t.Errorf("Expected message '%s', got '%s'", tt.expectedMsg, resp.Message)
			}

			if !resp.Error {
				t.Error("Expected error to be true")
			}
		})
	}
}
