package dto

import "customer-service/internal/core/domain"

type customerResponse struct {
	UserID      string `json:"id"`
	Name        string `json:"name"`
	ServiceArea int    `json:"service_area"`
}

func createCustomersResponse(customer domain.Customer) customerResponse {
	return customerResponse{
		UserID:      customer.UserID,
		Name:        customer.User.Name,
		ServiceArea: customer.ServiceArea,
	}
}

type CustomerListResponse []*customerResponse

func CreateCustomerListResponse(customers []domain.Customer) CustomerListResponse {
	response := CustomerListResponse{}
	for _, s := range customers {
		serviceArea := createCustomersResponse(s)
		response = append(response, &serviceArea)
	}
	return response
}
