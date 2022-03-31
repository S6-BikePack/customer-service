package ports

import (
	"customer-service/internal/core/domain"
)

type MessageBusPublisher interface {
	CreateCustomer(customer domain.Customer) error
	UpdateCustomerDetails(customer domain.Customer) error
	UpdateServiceArea(customer domain.Customer) error
}
