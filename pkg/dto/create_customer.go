package dto

import "customer-service/internal/core/domain"

type BodyCreateCustomer struct {
	Name     string
	LastName string
	Email    string
}

type ResponseCreateCustomer domain.Customer

func BuildResponseCreateCustomer(model domain.Customer) ResponseCreateCustomer {
	return ResponseCreateCustomer(model)
}
