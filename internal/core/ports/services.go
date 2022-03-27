package ports

import (
	"github.com/google/uuid"
	"rider-service/internal/core/domain"
)

type RiderService interface {
	GetAll() ([]domain.Rider, error)
	Get(uuid uuid.UUID) (domain.Rider, error)
	Create(name string, status int8) (domain.Rider, error)
	Update(uuid uuid.UUID, name string, status int8) (domain.Rider, error)
	UpdateLocation(uuid uuid.UUID, location domain.Location) (domain.Rider, error)
}
