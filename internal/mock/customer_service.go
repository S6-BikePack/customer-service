package mock

import (
	"context"
	"customer-service/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type CustomerService struct {
	mock.Mock
}

func (m *CustomerService) GetAll(ctx context.Context) ([]domain.Customer, error) {
	args := m.Called()
	return args.Get(0).([]domain.Customer), args.Error(1)
}

func (m *CustomerService) Get(ctx context.Context, id string) (domain.Customer, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Customer), args.Error(1)
}

func (m *CustomerService) Create(ctx context.Context, userId string, serviceArea int) (domain.Customer, error) {
	args := m.Called(userId, serviceArea)
	return args.Get(0).(domain.Customer), args.Error(1)
}

func (m *CustomerService) UpdateServiceArea(ctx context.Context, id string, serviceArea int) (domain.Customer, error) {
	args := m.Called(id, serviceArea)
	return args.Get(0).(domain.Customer), args.Error(1)
}

func (m *CustomerService) SaveOrUpdateUser(ctx context.Context, user domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}
