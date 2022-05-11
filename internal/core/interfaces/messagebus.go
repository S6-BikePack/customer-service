package interfaces

import (
	"context"
	"customer-service/internal/core/domain"
)

type MessageBusPublisher interface {
	CreateCustomer(ctx context.Context, customer domain.Customer) error
	UpdateServiceArea(ctx context.Context, customer domain.Customer) error
}
