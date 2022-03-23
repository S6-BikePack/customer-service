package domain

type Location struct {
	Latitude  float64
	Longitude float64
}

func NewLocation(latitude float64, longitude float64) Location {
	return Location{latitude, longitude}
}
