package mock

import (
	"context"
	"customer-service/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type CustomerRepository struct {
	mock.Mock
}

func (m *CustomerRepository) GetAll(ctx context.Context) ([]domain.Customer, error) {
	args := m.Called()
	return args.Get(0).([]domain.Customer), args.Error(1)
}

func (m *CustomerRepository) Get(ctx context.Context, id string) (domain.Customer, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Customer), args.Error(1)
}

func (m *CustomerRepository) Save(ctx context.Context, customer domain.Customer) (domain.Customer, error) {
	args := m.Called(customer)
	return args.Get(0).(domain.Customer), args.Error(1)
}

func (m *CustomerRepository) Update(ctx context.Context, customer domain.Customer) (domain.Customer, error) {
	args := m.Called(customer)
	return args.Get(0).(domain.Customer), args.Error(1)
}

func (m *CustomerRepository) SaveOrUpdateUser(ctx context.Context, user domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *CustomerRepository) GetUser(ctx context.Context, id string) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}
