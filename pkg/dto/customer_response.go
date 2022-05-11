package dto

import "customer-service/internal/core/domain"

type CustomerResponse struct {
	UserID      string `json:"id"`
	Name        string `json:"name"`
	LastName    string `json:"last_name"`
	ServiceArea int    `json:"service_area"`
}

func CreateCustomerResponse(customer domain.Customer) CustomerResponse {
	return CustomerResponse{
		UserID:      customer.UserID,
		Name:        customer.User.Name,
		LastName:    customer.User.LastName,
		ServiceArea: customer.ServiceArea,
	}
}
