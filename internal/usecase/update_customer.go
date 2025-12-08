package usecase

import (
	"context"
	"customer-service/internal/domain"
	"customer-service/internal/repository"
	"customer-service/pkg/errors"
)

type UpdateCustomerUseCase struct {
	repo repository.CustomerRepository
}

func NewUpdateCustomerUseCase(repo repository.CustomerRepository) *UpdateCustomerUseCase {
	return &UpdateCustomerUseCase{repo: repo}
}

func (uc *UpdateCustomerUseCase) Execute(ctx context.Context, id string, name, email *string) (*domain.Customer, error) {
	customer, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, errors.NewNotFoundError("Customer not found", "CUSTOMER_NOT_FOUND")
	}

	err = customer.Update(name, email)
	if err != nil {
		return nil, err
	}

	err = uc.repo.Update(ctx, customer)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
