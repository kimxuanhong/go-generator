package handler

// ErrorResponse represents a standard error payload returned by the API.
type ErrorResponse struct {
	Error     string `json:"error"`
	RequestID string `json:"request_id,omitempty"`
}
