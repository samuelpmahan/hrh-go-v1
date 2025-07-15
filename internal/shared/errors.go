package shared

import "fmt"

// AppError represents application-specific errors with codes and context
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common error codes
const (
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeInternal     = "INTERNAL_ERROR"
	ErrCodeDuplicate    = "DUPLICATE_ENTRY"
	ErrCodeConflict     = "CONFLICT"
)

// NewValidationError creates a validation error
func NewValidationError(message string, details ...string) *AppError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Code:    ErrCodeValidation,
		Message: message,
		Details: detail,
	}
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string, identifier string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Details: fmt.Sprintf("identifier: %s", identifier),
	}
}

// NewDuplicateError creates a duplicate entry error
func NewDuplicateError(resource string, field string, value string) *AppError {
	return &AppError{
		Code:    ErrCodeDuplicate,
		Message: fmt.Sprintf("%s already exists", resource),
		Details: fmt.Sprintf("%s: %s", field, value),
	}
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeUnauthorized,
		Message: message,
	}
}

// NewInternalError creates an internal error
func NewInternalError(message string, details ...string) *AppError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Code:    ErrCodeInternal,
		Message: message,
		Details: detail,
	}
}

// NewConflictError creates a conflict error
func NewConflictError(message string, details ...string) *AppError {
	var detail string
	if len(details) > 0 {
		detail = details[0]
	}
	return &AppError{
		Code:    ErrCodeConflict,
		Message: message,
		Details: detail,
	}
}
