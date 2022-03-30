package domain

import (
	"github.com/google/uuid"
)

type Rider struct {
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name     string
	Status   int8
	Location Location `gorm:"embedded"`
}

func NewRider(name string, status int8, location Location) Rider {
	return Rider{Name: name, Status: status, Location: location}
}
