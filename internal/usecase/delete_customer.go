package usecase

import (
	"context"
	"customer-service/internal/repository"
	"customer-service/pkg/errors"
	"fmt"
)

type DeleteCustomerUseCase struct {
	repo repository.CustomerRepository
}

func NewDeleteCustomerUseCase(repo repository.CustomerRepository) *DeleteCustomerUseCase {
	return &DeleteCustomerUseCase{repo: repo}
}

func (uc *DeleteCustomerUseCase) Execute(ctx context.Context, id string) error {
	customer, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if customer == nil {
		return errors.NewNotFoundError(
			fmt.Sprintf("Customer with id %s not found", id),
			"CUSTOMER_NOT_FOUND",
		)
	}

	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
