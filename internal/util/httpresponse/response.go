package httpresponse

import (
	"errors"
	"log"
	"net/http"

	"tms-core-service/internal/domain/errs"
	"tms-core-service/internal/util/apierror"

	"github.com/gofiber/fiber/v2"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Errors  interface{} `json:"errors,omitempty"` // For structured validation errors
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
	Meta    *Pagination `json:"meta"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Page   int   `json:"page"`
}

// Success sends a successful response
func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(http.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a created response (201)
func Created(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(http.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends a standardized error response
func Error(c *fiber.Ctx, err error) error {
	// Log the full error chain for internal debugging
	log.Printf("[ERROR] %v", err)

	var apiErr *apierror.APIError

	// Get Trace ID from context
	traceID, _ := c.Locals(apierror.CtxTraceID).(string)

	// Check if it's already an APIError
	if !errors.As(err, &apiErr) {
		// Map domain errors or generic errors to APIError
		apiErr = mapToAPIError(err)
	}

	// Always attach trace ID
	apiErr.TraceID = traceID

	return c.Status(apiErr.StatusCode).JSON(apiErr)
}

// Paginated sends a paginated response
func Paginated(c *fiber.Ctx, data interface{}, total int64, limit, offset int) error {
	page := (offset / limit) + 1
	if limit == 0 {
		page = 1
	}

	return c.Status(http.StatusOK).JSON(PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: &Pagination{
			Total:  total,
			Limit:  limit,
			Offset: offset,
			Page:   page,
		},
	})
}

// mapToAPIError maps various error types to a standardized APIError
func mapToAPIError(err error) *apierror.APIError {
	var valErrs errs.ValidationErrors
	if errors.As(err, &valErrs) {
		return apierror.NewValidationError("Validation failed", valErrs)
	}

	switch {
	case errors.Is(err, errs.ErrNotFound):
		return apierror.NewNotFoundError("Resource not found")
	case errors.Is(err, errs.ErrUnauthorized):
		return apierror.NewUnauthorizedError("Unauthorized access")
	case errors.Is(err, errs.ErrInvalidCredentials):
		return &apierror.APIError{
			Code:       apierror.CodeInvalidCredentials,
			Message:    "Invalid email or password",
			StatusCode: http.StatusUnauthorized,
		}
	case errors.Is(err, errs.ErrForbidden):
		return apierror.NewForbiddenError("Access forbidden")
	case errors.Is(err, errs.ErrConflict):
		return apierror.NewConflictError("Resource already exists")
	case errors.Is(err, errs.ErrBadRequest):
		return apierror.NewBadRequestError("Bad request parameters")
	case errors.Is(err, errs.ErrTokenExpired):
		return &apierror.APIError{
			Code:       apierror.CodeTokenExpired,
			Message:    "Authentication token has expired",
			StatusCode: http.StatusUnauthorized,
		}
	case errors.Is(err, errs.ErrTokenInvalid):
		return &apierror.APIError{
			Code:       apierror.CodeTokenInvalid,
			Message:    "Authentication token is invalid",
			StatusCode: http.StatusUnauthorized,
		}
	default:
		// Do not expose internal server errors
		return apierror.NewInternalError("")
	}
}
