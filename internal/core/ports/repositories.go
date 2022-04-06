package ports

import (
	"customer-service/internal/core/domain"
)

type CustomerRepository interface {
	GetAll() ([]domain.Customer, error)
	Get(id string) (domain.Customer, error)
	Save(customer domain.Customer) (domain.Customer, error)
	Update(customer domain.Customer) (domain.Customer, error)
	SaveOrUpdateCustomer(user domain.User) error
	GetUser(id string) (domain.User, error)
}
