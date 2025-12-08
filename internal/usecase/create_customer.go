package usecase

import (
	"context"
	"customer-service/internal/domain"
	"customer-service/internal/repository"
	"customer-service/pkg/errors"
)

type CreateCustomerUseCase struct {
	repo repository.CustomerRepository
}

func NewCreateCustomerUseCase(repo repository.CustomerRepository) *CreateCustomerUseCase {
	return &CreateCustomerUseCase{repo: repo}
}

func (uc *CreateCustomerUseCase) Execute(ctx context.Context, name, cpf, email string) (*domain.Customer, error) {
	customer, err := domain.NewCustomer(name, cpf, email)
	if err != nil {
		return nil, err
	}

	// Check if customer already exists
	existingCustomer, err := uc.repo.FindByCPFOrEmail(ctx, customer.CPF, customer.Email)
	if err != nil {
		return nil, err
	}
	if existingCustomer != nil {
		return nil, errors.NewConflictError("Customer already exists.", "CUSTOMER_ALREADY_EXISTS")
	}

	err = uc.repo.Create(ctx, customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
