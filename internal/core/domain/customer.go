package domain

import "github.com/google/uuid"

type Customer struct {
	ID          uuid.UUID `gorm:"type:uuid"`
	Name        string
	LastName    string
	Email       string
	ServiceArea int
}

func NewCustomer(name, lastName, email string, serviceArea int) Customer {
	return Customer{
		Name:        name,
		LastName:    lastName,
		Email:       email,
		ServiceArea: serviceArea,
	}
}
