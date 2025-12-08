package usecase

import (
	"context"
	"customer-service/internal/domain"
	"customer-service/internal/repository"
	"customer-service/pkg/errors"
	"customer-service/pkg/validator"
	"fmt"
)

type GetCustomerByCPFUseCase struct {
	repo repository.CustomerRepository
}

func NewGetCustomerByCPFUseCase(repo repository.CustomerRepository) *GetCustomerByCPFUseCase {
	return &GetCustomerByCPFUseCase{repo: repo}
}

func (uc *GetCustomerByCPFUseCase) Execute(ctx context.Context, cpf string) (*domain.Customer, error) {
	cleanCPF := validator.CleanCPF(cpf)
	if !validator.IsValidCPF(cleanCPF) {
		return nil, errors.NewValidationError("Invalid CPF", "INVALID_CPF")
	}

	customer, err := uc.repo.FindByCPF(ctx, cleanCPF)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, errors.NewNotFoundError(
			fmt.Sprintf("Customer with CPF %s not found", cpf),
			"CUSTOMER_NOT_FOUND",
		)
	}

	return customer, nil
}
