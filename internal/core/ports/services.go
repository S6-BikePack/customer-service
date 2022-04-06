package ports

import (
	"customer-service/internal/core/domain"
)

type CustomerService interface {
	GetAll() ([]domain.Customer, error)
	Get(id string) (domain.Customer, error)
	Create(userId string, serviceArea int) (domain.Customer, error)
	UpdateServiceArea(id string, serviceArea int) (domain.Customer, error)
	SaveOrUpdateUser(user domain.User) error
}
