package riderhdl

import "rider-service/internal/core/domain"

type BodyUpdate struct {
	Name   string
	Status string
}

type ResponseUpdate domain.Rider

func BuildResponseUpdate(model domain.Rider) ResponseCreate {
	return ResponseCreate(model)
}
