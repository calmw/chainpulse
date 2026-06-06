package chainpulse

import (
	"encoding/json"
	"fmt"
	"strings"
)

type APIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("chainpulse: status=%d message=%s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("chainpulse: status=%d body=%s", e.StatusCode, e.Body)
}

func newAPIError(statusCode int, body []byte) error {
	var payload struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	_ = json.Unmarshal(body, &payload)
	message := strings.TrimSpace(payload.Message)
	if message == "" {
		message = strings.TrimSpace(payload.Error)
	}
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Body:       string(body),
	}
}
