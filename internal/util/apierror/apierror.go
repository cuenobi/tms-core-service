package apierror

import (
	"fmt"
	"net/http"
)

// ErrorCode represents a machine-readable error code
type ErrorCode string

const (
	CodeValidationError    ErrorCode = "VALIDATION_ERROR"
	CodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	CodeForbidden          ErrorCode = "FORBIDDEN"
	CodeNotFound           ErrorCode = "NOT_FOUND"
	CodeConflict           ErrorCode = "CONFLICT"
	CodeInternalError      ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeBadRequest         ErrorCode = "BAD_REQUEST"
	CodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	CodeTokenExpired       ErrorCode = "TOKEN_EXPIRED"
	CodeTokenInvalid       ErrorCode = "TOKEN_INVALID"
)

const (
	HeaderXTraceID = "X-Trace-ID"
	CtxTraceID     = "traceId"
)

// APIError represents a standardized API error response
type APIError struct {
	Code       ErrorCode           `json:"code"`
	Message    string              `json:"message"`
	Errors     map[string][]string `json:"errors,omitempty"`
	TraceID    string              `json:"traceId,omitempty"`
	StatusCode int                 `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(message string, errors map[string][]string) *APIError {
	return &APIError{
		Code:       CodeValidationError,
		Message:    message,
		Errors:     errors,
		StatusCode: http.StatusBadRequest,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) *APIError {
	return &APIError{
		Code:       CodeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) *APIError {
	return &APIError{
		Code:       CodeForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) *APIError {
	return &APIError{
		Code:       CodeNotFound,
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) *APIError {
	return &APIError{
		Code:       CodeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// NewInternalError creates a new internal server error
func NewInternalError(message string) *APIError {
	if message == "" {
		message = "An unexpected error occurred"
	}
	return &APIError{
		Code:       CodeInternalError,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewBadRequestError creates a new bad request error
func NewBadRequestError(message string) *APIError {
	return &APIError{
		Code:       CodeBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}
