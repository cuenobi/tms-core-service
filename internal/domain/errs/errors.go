package errs

import "errors"

// Domain errors
var (
	// ErrNotFound indicates the requested resource was not found
	ErrNotFound = errors.New("resource not found")

	// ErrUnauthorized indicates authentication failure
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden indicates insufficient permissions
	ErrForbidden = errors.New("forbidden")

	// ErrConflict indicates a conflict with existing data (e.g., duplicate email)
	ErrConflict = errors.New("conflict")

	// ErrBadRequest indicates invalid input
	ErrBadRequest = errors.New("bad request")

	// ErrInternal indicates an internal server error
	ErrInternal = errors.New("internal server error")

	// ErrInvalidCredentials indicates invalid login credentials
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrTokenExpired indicates the JWT token has expired
	ErrTokenExpired = errors.New("token expired")

	// ErrTokenInvalid indicates the JWT token is invalid
	ErrTokenInvalid = errors.New("token invalid")
)

// ValidationError represents field-specific validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error returns the error message
func (e ValidationError) Error() string {
	return e.Message
}

// ValidationErrors is a collection of validation errors: field -> messages
type ValidationErrors map[string][]string

func (v ValidationErrors) Error() string {
	return "validation failed"
}
