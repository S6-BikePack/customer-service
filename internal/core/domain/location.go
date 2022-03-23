package domain

import "errors"

type Location struct {
	Latitude  float64
	Longitude float64
}

func NewLocation(latitude float64, longitude float64) (Location, error) {
	if latitude > 90 || latitude < -90 {
		return Location{}, errors.New("latitude is out of bounds")
	}

	if longitude > 180 || longitude < -180 {
		return Location{}, errors.New("longitude is out of bounds")
	}

	return Location{latitude, longitude}, nil
}
