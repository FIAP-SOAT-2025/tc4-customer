package usecase

import (
	"context"
	"customer-service/internal/domain"
	"customer-service/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateCustomerUseCase_Execute(t *testing.T) {
	newName := "Jane Doe"
	newEmail := "jane@example.com"
	invalidEmail := "invalid-email"

	tests := []struct {
		name          string
		customerID    string
		updateName    *string
		updateEmail   *string
		mockSetup     func(*MockCustomerRepository)
		expectError   bool
		expectedError string
	}{
		{
			name:        "Successfully update name and email",
			customerID:  "123",
			updateName:  &newName,
			updateEmail: &newEmail,
			mockSetup: func(m *MockCustomerRepository) {
				customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
				m.On("FindByID", mock.Anything, "123").
					Return(customer, nil)
				m.On("Update", mock.Anything, mock.Anything).
					Return(nil)
			},
			expectError: false,
		},
		{
			name:       "Successfully update only name",
			customerID: "123",
			updateName: &newName,
			mockSetup: func(m *MockCustomerRepository) {
				customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
				m.On("FindByID", mock.Anything, "123").
					Return(customer, nil)
				m.On("Update", mock.Anything, mock.Anything).
					Return(nil)
			},
			expectError: false,
		},
		{
			name:        "Successfully update only email",
			customerID:  "123",
			updateEmail: &newEmail,
			mockSetup: func(m *MockCustomerRepository) {
				customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
				m.On("FindByID", mock.Anything, "123").
					Return(customer, nil)
				m.On("Update", mock.Anything, mock.Anything).
					Return(nil)
			},
			expectError: false,
		},
		{
			name:       "Customer not found",
			customerID: "999",
			updateName: &newName,
			mockSetup: func(m *MockCustomerRepository) {
				m.On("FindByID", mock.Anything, "999").
					Return(nil, nil)
			},
			expectError:   true,
			expectedError: "CUSTOMER_NOT_FOUND",
		},
		{
			name:        "Invalid email",
			customerID:  "123",
			updateEmail: &invalidEmail,
			mockSetup: func(m *MockCustomerRepository) {
				customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
				m.On("FindByID", mock.Anything, "123").
					Return(customer, nil)
			},
			expectError:   true,
			expectedError: "INVALID_EMAIL",
		},
		{
			name:       "FindByID returns error",
			customerID: "123",
			updateName: &newName,
			mockSetup: func(m *MockCustomerRepository) {
				m.On("FindByID", mock.Anything, "123").
					Return(nil, errors.NewInternalError("database error"))
			},
			expectError: true,
		},
		{
			name:       "Update returns error",
			customerID: "123",
			updateName: &newName,
			mockSetup: func(m *MockCustomerRepository) {
				customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
				m.On("FindByID", mock.Anything, "123").
					Return(customer, nil)
				m.On("Update", mock.Anything, mock.Anything).
					Return(errors.NewInternalError("update failed"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCustomerRepository)
			tt.mockSetup(mockRepo)

			uc := NewUpdateCustomerUseCase(mockRepo)
			customer, err := uc.Execute(context.Background(), tt.customerID, tt.updateName, tt.updateEmail)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, customer)
				if tt.expectedError != "" {
					appErr, ok := err.(*errors.AppError)
					assert.True(t, ok)
					assert.Equal(t, tt.expectedError, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, customer)
				if tt.updateName != nil {
					assert.Equal(t, *tt.updateName, customer.Name)
				}
				if tt.updateEmail != nil {
					assert.Equal(t, *tt.updateEmail, customer.Email)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
