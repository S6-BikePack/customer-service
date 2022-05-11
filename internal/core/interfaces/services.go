package interfaces

import (
	"context"
	"customer-service/internal/core/domain"
)

type CustomerService interface {
	GetAll(ctx context.Context) ([]domain.Customer, error)
	Get(ctx context.Context, id string) (domain.Customer, error)
	Create(ctx context.Context, userId string, serviceArea int) (domain.Customer, error)
	UpdateServiceArea(ctx context.Context, id string, serviceArea int) (domain.Customer, error)
	SaveOrUpdateUser(ctx context.Context, user domain.User) error
}
