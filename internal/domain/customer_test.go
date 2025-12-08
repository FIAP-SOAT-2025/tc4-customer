package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCustomer(t *testing.T) {
	tests := []struct {
		name        string
		customerName string
		cpf         string
		email       string
		expectError bool
		errorCode   string
	}{
		{
			name:         "Valid customer",
			customerName: "John Doe",
			cpf:          "111.444.777-35",
			email:        "john@example.com",
			expectError:  false,
		},
		{
			name:         "Empty name",
			customerName: "",
			cpf:          "11144477735",
			email:        "john@example.com",
			expectError:  true,
			errorCode:    "NAME_EMPTY",
		},
		{
			name:         "Invalid CPF",
			customerName: "John Doe",
			cpf:          "12345678901",
			email:        "john@example.com",
			expectError:  true,
			errorCode:    "INVALID_CPF",
		},
		{
			name:         "Invalid email",
			customerName: "John Doe",
			cpf:          "11144477735",
			email:        "invalid-email",
			expectError:  true,
			errorCode:    "INVALID_EMAIL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer, err := NewCustomer(tt.customerName, tt.cpf, tt.email)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, customer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, customer)
				assert.Equal(t, tt.customerName, customer.Name)
				assert.NotEmpty(t, customer.ID)
				assert.NotZero(t, customer.CreatedAt)
				assert.NotZero(t, customer.UpdatedAt)
			}
		})
	}
}

func TestCustomer_Update(t *testing.T) {
	customer, _ := NewCustomer("John Doe", "11144477735", "john@example.com")
	initialUpdatedAt := customer.UpdatedAt

	tests := []struct {
		name        string
		newName     *string
		newEmail    *string
		expectError bool
	}{
		{
			name:        "Update name only",
			newName:     stringPtr("Jane Doe"),
			newEmail:    nil,
			expectError: false,
		},
		{
			name:        "Update email only",
			newName:     nil,
			newEmail:    stringPtr("jane@example.com"),
			expectError: false,
		},
		{
			name:        "Update both",
			newName:     stringPtr("Jane Smith"),
			newEmail:    stringPtr("jane.smith@example.com"),
			expectError: false,
		},
		{
			name:        "Invalid empty name",
			newName:     stringPtr(""),
			newEmail:    nil,
			expectError: true,
		},
		{
			name:        "Invalid email",
			newName:     nil,
			newEmail:    stringPtr("invalid"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := customer.Update(tt.newName, tt.newEmail)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, customer.UpdatedAt.After(initialUpdatedAt))

				if tt.newName != nil {
					assert.Equal(t, *tt.newName, customer.Name)
				}
				if tt.newEmail != nil {
					assert.Equal(t, *tt.newEmail, customer.Email)
				}
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
