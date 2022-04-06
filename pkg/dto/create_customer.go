package dto

import "customer-service/internal/core/domain"

type BodyCreateCustomer struct {
	ID          string
	ServiceArea int
}

type ResponseCreateCustomer domain.Customer

func BuildResponseCreateCustomer(model domain.Customer) ResponseCreateCustomer {
	return ResponseCreateCustomer(model)
}
