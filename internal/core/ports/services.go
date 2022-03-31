package ports

import (
	"customer-service/internal/core/domain"
	"github.com/google/uuid"
)

type CustomerService interface {
	GetAll() ([]domain.Customer, error)
	Get(uuid uuid.UUID) (domain.Customer, error)
	Create(name, lastName, email string) (domain.Customer, error)
	UpdateCustomerDetails(uuid uuid.UUID, name, lastName, email string) (domain.Customer, error)
	UpdateServiceArea(uuid uuid.UUID, serviceArea int) (domain.Customer, error)
}
