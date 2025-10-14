package errs

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode represents a machine-readable error identifier surfaced to clients.
type ErrorCode string

const (
	CodeInternal     ErrorCode = "internal_error"
	CodeNotFound     ErrorCode = "not_found"
	CodeInvalidInput ErrorCode = "invalid_input"
	CodeUnauthorized ErrorCode = "unauthorized"
	CodeConflict     ErrorCode = "conflict"
)

// AppError keeps domain error details together with HTTP semantics.
type AppError struct {
	Code    ErrorCode
	Message string
	Status  int
	Err     error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// New creates a new AppError with the provided parameters.
func New(code ErrorCode, status int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Status:  status,
		Message: message,
		Err:     err,
	}
}

// From converts an arbitrary error into an AppError, defaulting to an internal server error.
func From(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return New(CodeInternal, http.StatusInternalServerError, "internal server error", err)
}
