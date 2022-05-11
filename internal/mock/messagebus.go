package mock

import (
	"context"
	"customer-service/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MessageBusPublisher struct {
	mock.Mock
}

func (m *MessageBusPublisher) CreateCustomer(ctx context.Context, customer domain.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *MessageBusPublisher) UpdateServiceArea(ctx context.Context, customer domain.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}
