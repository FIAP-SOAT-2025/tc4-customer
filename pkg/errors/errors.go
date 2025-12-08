package errors

import "fmt"

type AppError struct {
	Message    string
	StatusCode int
	Code       string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewValidationError(message, code string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: 400,
		Code:       code,
	}
}

func NewNotFoundError(message, code string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: 404,
		Code:       code,
	}
}

func NewConflictError(message, code string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: 409,
		Code:       code,
	}
}

func NewInternalError(message string) *AppError {
	return &AppError{
		Message:    message,
		StatusCode: 500,
		Code:       "INTERNAL_ERROR",
	}
}

func WrapError(err error, message string) *AppError {
	return &AppError{
		Message:    fmt.Sprintf("%s: %v", message, err),
		StatusCode: 500,
		Code:       "INTERNAL_ERROR",
	}
}
