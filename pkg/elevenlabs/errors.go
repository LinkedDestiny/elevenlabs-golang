package elevenlabs

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ElevenLabsError represents an error from the ElevenLabs API
type ElevenLabsError interface {
	error
	StatusCode() int
	Body() interface{}
}

// APIError represents a generic API error
type APIError struct {
	statusCode int
	body       interface{}
	message    string
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.message != "" {
		return e.message
	}
	return fmt.Sprintf("API error with status code %d", e.statusCode)
}

// StatusCode returns the HTTP status code
func (e *APIError) StatusCode() int {
	return e.statusCode
}

// Body returns the error response body
func (e *APIError) Body() interface{} {
	return e.body
}

// BadRequestError represents a 400 Bad Request error
type BadRequestError struct {
	*APIError
}

// NewBadRequestError creates a new BadRequestError
func NewBadRequestError(body interface{}) *BadRequestError {
	return &BadRequestError{
		APIError: &APIError{
			statusCode: 400,
			body:       body,
			message:    "Bad Request",
		},
	}
}

// ForbiddenError represents a 403 Forbidden error
type ForbiddenError struct {
	*APIError
}

// NewForbiddenError creates a new ForbiddenError
func NewForbiddenError(body interface{}) *ForbiddenError {
	return &ForbiddenError{
		APIError: &APIError{
			statusCode: 403,
			body:       body,
			message:    "Forbidden",
		},
	}
}

// NotFoundError represents a 404 Not Found error
type NotFoundError struct {
	*APIError
}

// NewNotFoundError creates a new NotFoundError
func NewNotFoundError(body interface{}) *NotFoundError {
	return &NotFoundError{
		APIError: &APIError{
			statusCode: 404,
			body:       body,
			message:    "Not Found",
		},
	}
}

// UnprocessableEntityError represents a 422 Unprocessable Entity error
type UnprocessableEntityError struct {
	*APIError
}

// NewUnprocessableEntityError creates a new UnprocessableEntityError
func NewUnprocessableEntityError(body interface{}) *UnprocessableEntityError {
	return &UnprocessableEntityError{
		APIError: &APIError{
			statusCode: 422,
			body:       body,
			message:    "Unprocessable Entity",
		},
	}
}

// ErrorResponse represents a structured error response from the API
type ErrorResponse struct {
	Detail  interface{} `json:"detail,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ParseAPIError parses an HTTP response and returns an appropriate error
func ParseAPIError(resp *http.Response) error {
	if resp.StatusCode < 400 {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIError{
			statusCode: resp.StatusCode,
			message:    fmt.Sprintf("Failed to read error response: %v", err),
		}
	}

	var errorResp ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		// If we can't parse as JSON, use the raw body
		return &APIError{
			statusCode: resp.StatusCode,
			body:       string(body),
			message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)),
		}
	}

	switch resp.StatusCode {
	case 400:
		return NewBadRequestError(errorResp)
	case 403:
		return NewForbiddenError(errorResp)
	case 404:
		return NewNotFoundError(errorResp)
	case 422:
		return NewUnprocessableEntityError(errorResp)
	default:
		message := errorResp.Error
		if message == "" {
			message = errorResp.Message
		}
		if message == "" {
			message = fmt.Sprintf("HTTP %d error", resp.StatusCode)
		}
		return &APIError{
			statusCode: resp.StatusCode,
			body:       errorResp,
			message:    message,
		}
	}
}
