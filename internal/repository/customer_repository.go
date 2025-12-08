package repository

import (
	"context"
	"customer-service/internal/domain"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	FindByID(ctx context.Context, id string) (*domain.Customer, error)
	FindByCPF(ctx context.Context, cpf string) (*domain.Customer, error)
	FindByCPFOrEmail(ctx context.Context, cpf, email string) (*domain.Customer, error)
	Update(ctx context.Context, customer *domain.Customer) error
	Delete(ctx context.Context, id string) error
	GetEmailByID(ctx context.Context, id string) (string, error)
}
