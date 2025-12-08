package domain

import (
	"customer-service/pkg/errors"
	"customer-service/pkg/validator"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	CPF       string    `json:"cpf" bson:"cpf"`
	Email     string    `json:"email" bson:"email"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func NewCustomer(name, cpf, email string) (*Customer, error) {
	// Validate name
	if strings.TrimSpace(name) == "" {
		return nil, errors.NewValidationError("Name cannot be empty", "NAME_EMPTY")
	}

	// Validate and clean CPF
	cleanCPF := validator.CleanCPF(cpf)
	if !validator.IsValidCPF(cleanCPF) {
		return nil, errors.NewValidationError("Invalid CPF", "INVALID_CPF")
	}

	// Validate email
	cleanEmail := strings.ToLower(strings.TrimSpace(email))
	if !validator.IsValidEmail(cleanEmail) {
		return nil, errors.NewValidationError("Invalid Email", "INVALID_EMAIL")
	}

	now := time.Now()
	return &Customer{
		ID:        uuid.New().String(),
		Name:      name,
		CPF:       cleanCPF,
		Email:     cleanEmail,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (c *Customer) Update(name, email *string) error {
	if name != nil {
		if strings.TrimSpace(*name) == "" {
			return errors.NewValidationError("Name cannot be empty", "NAME_EMPTY")
		}
		c.Name = *name
	}

	if email != nil {
		cleanEmail := strings.ToLower(strings.TrimSpace(*email))
		if !validator.IsValidEmail(cleanEmail) {
			return errors.NewValidationError("Invalid Email", "INVALID_EMAIL")
		}
		c.Email = cleanEmail
	}

	c.UpdatedAt = time.Now()
	return nil
}
