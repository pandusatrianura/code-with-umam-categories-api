package json_wrapper

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name            string
		requestBody     string
		expectedError   string // Leave empty if no error is expected
		payload         any
		expectedPayload any
	}{
		{
			name:            "valid JSON into map",
			requestBody:     `{"key":"value"}`,
			expectedError:   "",
			payload:         &map[string]string{},
			expectedPayload: &map[string]string{"key": "value"},
		},
		{
			name:            "valid JSON into struct",
			requestBody:     `{"key":"value"}`,
			expectedError:   "",
			payload:         &struct{ Key string }{},
			expectedPayload: &struct{ Key string }{Key: "value"},
		},
		{
			name:          "malformed JSON",
			requestBody:   `{"key":`,
			expectedError: "unexpected EOF",
			payload:       &map[string]string{},
		},
		{
			name:            "empty JSON object",
			requestBody:     `{}`,
			expectedError:   "",
			payload:         &map[string]string{},
			expectedPayload: &map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.requestBody)))

			err := ParseJSON(r, tt.payload)

			// Check error
			if err != nil {
				if tt.expectedError == "" || err.Error() != tt.expectedError {
					t.Errorf("unexpected error: got %v, want %v", err, tt.expectedError)
				}
			} else if tt.expectedError != "" {
				t.Errorf("expected error but got none, want: %v", tt.expectedError)
			}

			// Check payload if no error is expected
			if tt.expectedError == "" && tt.payload != nil {
				if fmt.Sprintf("%v", tt.payload) != fmt.Sprintf("%v", tt.expectedPayload) {
					t.Errorf("unexpected payload: got %v, want %v", tt.payload, tt.expectedPayload)
				}
			}
		})
	}
}

func TestWriteJSONResponse(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		input        any
		expectedBody string
		expectedType string
	}{
		{
			name:         "valid struct response",
			status:       http.StatusOK,
			input:        map[string]string{"key": "value"},
			expectedBody: `{"key":"value"}`,
			expectedType: "application/json",
		},
		{
			name:         "response with nil body",
			status:       http.StatusNoContent,
			input:        nil,
			expectedBody: "null",
			expectedType: "application/json",
		},
		{
			name:         "response with slice",
			status:       http.StatusOK,
			input:        []string{"item1", "item2"},
			expectedBody: `["item1","item2"]`,
			expectedType: "application/json",
		},
		{
			name:         "response with error status",
			status:       http.StatusInternalServerError,
			input:        map[string]string{"error": "something went wrong"},
			expectedBody: `{"error":"something went wrong"}`,
			expectedType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			WriteJSONResponse(rec, tt.status, tt.input)

			// Check Status Code
			if rec.Code != tt.status {
				t.Errorf("incorrect status code, got: %d, want: %d", rec.Code, tt.status)
			}

			// Check Content-Type
			if contentType := rec.Header().Get("Content-Type"); contentType != tt.expectedType {
				t.Errorf("incorrect content type, got: %s, want: %s", contentType, tt.expectedType)
			}

			// Check Body
			body := bytes.TrimSpace(rec.Body.Bytes())
			if string(body) != tt.expectedBody {
				t.Errorf("incorrect body, got: %s, want: %s", body, tt.expectedBody)
			}
		})
	}
}
