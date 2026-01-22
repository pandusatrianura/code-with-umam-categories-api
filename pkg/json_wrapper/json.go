package json_wrapper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIResponse represents the standardized structure for API responses, encapsulating status, message, and optional data.
type APIResponse struct {
	Code    string      `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteJSONResponse writes a JSON response with the provided status code and value to the http.ResponseWriter.
func WriteJSONResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// ParseJSON decodes a JSON body from an HTTP request into the specified payload object.
func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing body request")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}
