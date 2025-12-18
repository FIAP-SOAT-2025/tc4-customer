package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	err := &AppError{
		Message:    "Test error message",
		StatusCode: 400,
		Code:       "TEST_ERROR",
	}

	assert.Equal(t, "Test error message", err.Error())
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("Invalid input", "INVALID_INPUT")

	assert.NotNil(t, err)
	assert.Equal(t, "Invalid input", err.Message)
	assert.Equal(t, 400, err.StatusCode)
	assert.Equal(t, "INVALID_INPUT", err.Code)
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("Resource not found", "NOT_FOUND")

	assert.NotNil(t, err)
	assert.Equal(t, "Resource not found", err.Message)
	assert.Equal(t, 404, err.StatusCode)
	assert.Equal(t, "NOT_FOUND", err.Code)
}

func TestNewConflictError(t *testing.T) {
	err := NewConflictError("Resource already exists", "ALREADY_EXISTS")

	assert.NotNil(t, err)
	assert.Equal(t, "Resource already exists", err.Message)
	assert.Equal(t, 409, err.StatusCode)
	assert.Equal(t, "ALREADY_EXISTS", err.Code)
}

func TestNewInternalError(t *testing.T) {
	err := NewInternalError("Internal server error")

	assert.NotNil(t, err)
	assert.Equal(t, "Internal server error", err.Message)
	assert.Equal(t, 500, err.StatusCode)
	assert.Equal(t, "INTERNAL_ERROR", err.Code)
}

func TestWrapError(t *testing.T) {
	originalErr := fmt.Errorf("original error")
	wrappedErr := WrapError(originalErr, "Failed to process")

	assert.NotNil(t, wrappedErr)
	assert.Contains(t, wrappedErr.Message, "Failed to process")
	assert.Contains(t, wrappedErr.Message, "original error")
	assert.Equal(t, 500, wrappedErr.StatusCode)
	assert.Equal(t, "INTERNAL_ERROR", wrappedErr.Code)
}
