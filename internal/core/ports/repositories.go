package ports

import (
	"customer-service/internal/core/domain"
	"github.com/google/uuid"
)

type CustomerRepository interface {
	GetAll() ([]domain.Customer, error)
	Get(uuid uuid.UUID) (domain.Customer, error)
	Save(customer domain.Customer) (domain.Customer, error)
	Update(customer domain.Customer) (domain.Customer, error)
}
