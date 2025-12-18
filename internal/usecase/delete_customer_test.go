package usecase

import (
	"context"
	"customer-service/internal/domain"
	"customer-service/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCustomerUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		customerID    string
		mockSetup     func(*MockCustomerRepository)
		expectError   bool
		expectedError string
	}{
		{
			name:       "Successfully delete customer",
			customerID: "123",
			mockSetup: func(m *MockCustomerRepository) {
				customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
				m.On("FindByID", mock.Anything, "123").
					Return(customer, nil)
				m.On("Delete", mock.Anything, "123").
					Return(nil)
			},
			expectError: false,
		},
		{
			name:       "Customer not found",
			customerID: "999",
			mockSetup: func(m *MockCustomerRepository) {
				m.On("FindByID", mock.Anything, "999").
					Return(nil, nil)
			},
			expectError:   true,
			expectedError: "CUSTOMER_NOT_FOUND",
		},
		{
			name:       "FindByID returns error",
			customerID: "123",
			mockSetup: func(m *MockCustomerRepository) {
				m.On("FindByID", mock.Anything, "123").
					Return(nil, errors.NewInternalError("database error"))
			},
			expectError: true,
		},
		{
			name:       "Delete returns error",
			customerID: "123",
			mockSetup: func(m *MockCustomerRepository) {
				customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
				m.On("FindByID", mock.Anything, "123").
					Return(customer, nil)
				m.On("Delete", mock.Anything, "123").
					Return(errors.NewInternalError("delete failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCustomerRepository)
			tt.mockSetup(mockRepo)

			uc := NewDeleteCustomerUseCase(mockRepo)
			err := uc.Execute(context.Background(), tt.customerID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != "" {
					appErr, ok := err.(*errors.AppError)
					assert.True(t, ok)
					assert.Equal(t, tt.expectedError, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
