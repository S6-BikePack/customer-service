package riderhdl

import "rider-service/internal/core/domain"

type BodyCreate struct {
	Name   string
	Status string
}

type ResponseCreate domain.Rider

func BuildResponseCreate(model domain.Rider) ResponseCreate {
	return ResponseCreate(model)
}
