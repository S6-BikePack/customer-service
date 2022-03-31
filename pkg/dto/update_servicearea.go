package dto

import "customer-service/internal/core/domain"

type BodyUpdateServiceArea struct {
	ServiceArea int
}

type ResponseUpdateServiceArea domain.Customer

func BuildResponseUpdateServiceArea(model domain.Customer) ResponseUpdateServiceArea {
	return ResponseUpdateServiceArea(model)
}
