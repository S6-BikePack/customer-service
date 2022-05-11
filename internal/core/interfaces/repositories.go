package interfaces

import (
	"context"
	"customer-service/internal/core/domain"
)

type CustomerRepository interface {
	GetAll(ctx context.Context) ([]domain.Customer, error)
	Get(ctx context.Context, id string) (domain.Customer, error)
	Save(ctx context.Context, customer domain.Customer) (domain.Customer, error)
	Update(ctx context.Context, customer domain.Customer) (domain.Customer, error)
	SaveOrUpdateUser(ctx context.Context, user domain.User) error
	GetUser(ctx context.Context, id string) (domain.User, error)
}
