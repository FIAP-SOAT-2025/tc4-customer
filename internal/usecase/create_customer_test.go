package usecase

import (
	"context"
	"customer-service/internal/domain"
	"customer-service/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCustomerRepository struct {
	mock.Mock
}

func (m *MockCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) FindByID(ctx context.Context, id string) (*domain.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockCustomerRepository) FindByCPF(ctx context.Context, cpf string) (*domain.Customer, error) {
	args := m.Called(ctx, cpf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockCustomerRepository) FindByCPFOrEmail(ctx context.Context, cpf, email string) (*domain.Customer, error) {
	args := m.Called(ctx, cpf, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockCustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockCustomerRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCustomerRepository) GetEmailByID(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func TestCreateCustomerUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		customerName  string
		cpf           string
		email         string
		mockSetup     func(*MockCustomerRepository)
		expectError   bool
		expectedError string
	}{
		{
			name:         "Successfully create customer",
			customerName: "John Doe",
			cpf:          "11144477735",
			email:        "john@example.com",
			mockSetup: func(m *MockCustomerRepository) {
				m.On("FindByCPFOrEmail", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, nil)
				m.On("Create", mock.Anything, mock.Anything).
					Return(nil)
			},
			expectError: false,
		},
		{
			name:         "Customer already exists",
			customerName: "John Doe",
			cpf:          "11144477735",
			email:        "john@example.com",
			mockSetup: func(m *MockCustomerRepository) {
				existingCustomer, _ := domain.NewCustomer("Jane Doe", "11144477735", "jane@example.com")
				m.On("FindByCPFOrEmail", mock.Anything, mock.Anything, mock.Anything).
					Return(existingCustomer, nil)
			},
			expectError:   true,
			expectedError: "CUSTOMER_ALREADY_EXISTS",
		},
		{
			name:         "Invalid CPF",
			customerName: "John Doe",
			cpf:          "invalid",
			email:        "john@example.com",
			mockSetup:    func(m *MockCustomerRepository) {},
			expectError:  true,
		},
		{
			name:         "Invalid email",
			customerName: "John Doe",
			cpf:          "11144477735",
			email:        "invalid",
			mockSetup:    func(m *MockCustomerRepository) {},
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockCustomerRepository)
			tt.mockSetup(mockRepo)

			uc := NewCreateCustomerUseCase(mockRepo)
			customer, err := uc.Execute(context.Background(), tt.customerName, tt.cpf, tt.email)

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
