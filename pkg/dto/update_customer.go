package dto

import "customer-service/internal/core/domain"

type BodyUpdateCustomer struct {
	Name     string
	LastName string
	Email    string
}

type ResponseUpdateCustomer domain.Customer

func BuildResponseUpdateCustomer(model domain.Customer) ResponseUpdateCustomer {
	return ResponseUpdateCustomer(model)
}
