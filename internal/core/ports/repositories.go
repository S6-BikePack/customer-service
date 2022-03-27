package ports

import (
	"github.com/google/uuid"
	"rider-service/internal/core/domain"
)

type RiderRepository interface {
	GetAll() ([]domain.Rider, error)
	Get(uuid uuid.UUID) (domain.Rider, error)
	Save(rider domain.Rider) (domain.Rider, error)
	Update(rider domain.Rider) (domain.Rider, error)
}
