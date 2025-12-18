package handler

import (
	"bytes"
	"context"
	"customer-service/internal/domain"
	"customer-service/internal/usecase"
	"customer-service/pkg/errors"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository for testing handlers
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, customer *domain.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*domain.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockRepository) FindByCPF(ctx context.Context, cpf string) (*domain.Customer, error) {
	args := m.Called(ctx, cpf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockRepository) FindByCPFOrEmail(ctx context.Context, cpf, email string) (*domain.Customer, error) {
	args := m.Called(ctx, cpf, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Customer), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, customer *domain.Customer) error {
	args := m.Called(ctx, customer)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetEmailByID(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func setupTestRouter(handler *CustomerHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/customer", handler.CreateCustomer)
	router.GET("/customer/:cpf", handler.GetCustomerByCPF)
	router.PATCH("/customer/:id", handler.UpdateCustomer)
	router.DELETE("/customer/:id", handler.DeleteCustomer)

	return router
}

func TestCreateCustomer(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockRepository)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successfully create customer",
			requestBody: CreateCustomerRequest{
				Name:  "John Doe",
				CPF:   "111.444.777-35",
				Email: "john@example.com",
			},
			mockSetup: func(m *MockRepository) {
				m.On("FindByCPFOrEmail", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, nil)
				m.On("Create", mock.Anything, mock.Anything).
					Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid request body",
			requestBody:    "invalid json",
			mockSetup:      func(m *MockRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
		},
		{
			name: "Missing required field",
			requestBody: map[string]string{
				"name": "John Doe",
			},
			mockSetup:      func(m *MockRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
		},
		{
			name: "Customer already exists",
			requestBody: CreateCustomerRequest{
				Name:  "John Doe",
				CPF:   "111.444.777-35",
				Email: "john@example.com",
			},
			mockSetup: func(m *MockRepository) {
				existing, _ := domain.NewCustomer("Jane", "11144477735", "jane@example.com")
				m.On("FindByCPFOrEmail", mock.Anything, mock.Anything, mock.Anything).
					Return(existing, nil)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "CUSTOMER_ALREADY_EXISTS",
		},
		{
			name: "Invalid CPF",
			requestBody: CreateCustomerRequest{
				Name:  "John Doe",
				CPF:   "invalid",
				Email: "john@example.com",
			},
			mockSetup:      func(m *MockRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_CPF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			handler := NewCustomerHandler(
				usecase.NewCreateCustomerUseCase(mockRepo),
				usecase.NewGetCustomerByCPFUseCase(mockRepo),
				usecase.NewUpdateCustomerUseCase(mockRepo),
				usecase.NewDeleteCustomerUseCase(mockRepo),
			)
			router := setupTestRouter(handler)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/customer", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedError, response["error"])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetCustomerByCPF(t *testing.T) {
	tests := []struct {
		name           string
		cpf            string
		mockSetup      func(*MockRepository)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successfully get customer",
			cpf:  "11144477735",
			mockSetup: func(m *MockRepository) {
				customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
				m.On("FindByCPF", mock.Anything, "11144477735").
					Return(customer, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Customer not found",
			cpf:  "11144477735",
			mockSetup: func(m *MockRepository) {
				m.On("FindByCPF", mock.Anything, "11144477735").
					Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "CUSTOMER_NOT_FOUND",
		},
		{
			name:           "Invalid CPF",
			cpf:            "invalid",
			mockSetup:      func(m *MockRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_CPF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			handler := NewCustomerHandler(
				usecase.NewCreateCustomerUseCase(mockRepo),
				usecase.NewGetCustomerByCPFUseCase(mockRepo),
				usecase.NewUpdateCustomerUseCase(mockRepo),
				usecase.NewDeleteCustomerUseCase(mockRepo),
			)
			router := setupTestRouter(handler)

			req := httptest.NewRequest(http.MethodGet, "/customer/"+tt.cpf, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedError, response["error"])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateCustomer(t *testing.T) {
	name := "Jane Doe"
	email := "jane@example.com"

	tests := []struct {
		name           string
		customerID     string
		requestBody    interface{}
		mockSetup      func(*MockRepository)
		expectedStatus int
		expectedError  string
	}{
		{
			name:       "Successfully update customer",
			customerID: "123",
			requestBody: UpdateCustomerRequest{
				Name:  &name,
				Email: &email,
			},
			mockSetup: func(m *MockRepository) {
				customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
				m.On("FindByID", mock.Anything, "123").
					Return(customer, nil)
				m.On("Update", mock.Anything, mock.Anything).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid request body",
			customerID:     "123",
			requestBody:    "invalid json",
			mockSetup:      func(m *MockRepository) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
		},
		{
			name:       "Customer not found",
			customerID: "999",
			requestBody: UpdateCustomerRequest{
				Name: &name,
			},
			mockSetup: func(m *MockRepository) {
				m.On("FindByID", mock.Anything, "999").
					Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "CUSTOMER_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			handler := NewCustomerHandler(
				usecase.NewCreateCustomerUseCase(mockRepo),
				usecase.NewGetCustomerByCPFUseCase(mockRepo),
				usecase.NewUpdateCustomerUseCase(mockRepo),
				usecase.NewDeleteCustomerUseCase(mockRepo),
			)
			router := setupTestRouter(handler)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/customer/"+tt.customerID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedError, response["error"])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteCustomer(t *testing.T) {
	tests := []struct {
		name           string
		customerID     string
		mockSetup      func(*MockRepository)
		expectedStatus int
		expectedError  string
	}{
		{
			name:       "Successfully delete customer",
			customerID: "123",
			mockSetup: func(m *MockRepository) {
				customer, _ := domain.NewCustomer("John Doe", "11144477735", "john@example.com")
				m.On("FindByID", mock.Anything, "123").
					Return(customer, nil)
				m.On("Delete", mock.Anything, "123").
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:       "Customer not found",
			customerID: "999",
			mockSetup: func(m *MockRepository) {
				m.On("FindByID", mock.Anything, "999").
					Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "CUSTOMER_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.mockSetup(mockRepo)

			handler := NewCustomerHandler(
				usecase.NewCreateCustomerUseCase(mockRepo),
				usecase.NewGetCustomerByCPFUseCase(mockRepo),
				usecase.NewUpdateCustomerUseCase(mockRepo),
				usecase.NewDeleteCustomerUseCase(mockRepo),
			)
			router := setupTestRouter(handler)

			req := httptest.NewRequest(http.MethodDelete, "/customer/"+tt.customerID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Equal(t, tt.expectedError, response["error"])
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandleError(t *testing.T) {
	tests := []struct {
		name           string
		error          error
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "AppError - Validation",
			error:          errors.NewValidationError("Invalid input", "INVALID_INPUT"),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_INPUT",
		},
		{
			name:           "AppError - Not Found",
			error:          errors.NewNotFoundError("Not found", "NOT_FOUND"),
			expectedStatus: http.StatusNotFound,
			expectedError:  "NOT_FOUND",
		},
		{
			name:           "AppError - Conflict",
			error:          errors.NewConflictError("Conflict", "CONFLICT"),
			expectedStatus: http.StatusConflict,
			expectedError:  "CONFLICT",
		},
		{
			name:           "Generic Error",
			error:          assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handleError(c, tt.error)

			assert.Equal(t, tt.expectedStatus, w.Code)
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedError, response["error"])
		})
	}
}
