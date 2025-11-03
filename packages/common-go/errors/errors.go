package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	ErrorCodeInvalidInput   ErrorCode = "INVALID_INPUT"
	ErrorCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrorCodeConflict       ErrorCode = "CONFLICT"
	ErrorCodeInternal       ErrorCode = "INTERNAL_ERROR"
	ErrorCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

// AppError represents an application error
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	HTTPStatus int       `json:"-"`
	Details    string    `json:"details,omitempty"`
	Err        error     `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAppError creates a new application error
func NewAppError(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// NewInvalidInputError creates an invalid input error
func NewInvalidInputError(message string) *AppError {
	return NewAppError(ErrorCodeInvalidInput, message, http.StatusBadRequest)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return NewAppError(ErrorCodeUnauthorized, message, http.StatusUnauthorized)
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string) *AppError {
	return NewAppError(ErrorCodeForbidden, message, http.StatusForbidden)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(message string) *AppError {
	return NewAppError(ErrorCodeNotFound, message, http.StatusNotFound)
}

// NewConflictError creates a conflict error
func NewConflictError(message string) *AppError {
	return NewAppError(ErrorCodeConflict, message, http.StatusConflict)
}

// NewInternalError creates an internal error
func NewInternalError(message string) *AppError {
	return NewAppError(ErrorCodeInternal, message, http.StatusInternalServerError)
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// WithError wraps an underlying error
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

