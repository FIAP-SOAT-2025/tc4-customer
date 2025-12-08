package usecase

import (
	"context"
	"customer-service/internal/domain"
	"customer-service/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCustomerByCPFUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		cpf           string
		mockSetup     func(*MockCustomerRepository)
		expectError   bool
		expectedError string
	}{
		{
			name: "Successfully get customer by CPF",
			cpf:  "11144477735",
			mockSetup: func(m *MockCustomerRepository) {
				customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
				m.On("FindByCPF", mock.Anything, "11144477735").
					Return(customer, nil)
			},
			expectError: false,
		},
		{
			name: "Customer not found",
			cpf:  "11144477735",
			mockSetup: func(m *MockCustomerRepository) {
				m.On("FindByCPF", mock.Anything, "11144477735").
					Return(nil, nil)
			},
			expectError:   true,
			expectedError: "CUSTOMER_NOT_FOUND",
		},
		{
			name:          "Invalid CPF",
			cpf:           "invalid",
			mockSetup:     func(m *MockCustomerRepository) {},
			expectError:   true,
			expectedError: "INVALID_CPF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCustomerRepository)
			tt.mockSetup(mockRepo)

			uc := NewGetCustomerByCPFUseCase(mockRepo)
			customer, err := uc.Execute(context.Background(), tt.cpf)

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
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
