package ports

import "rider-service/internal/core/domain"

type RiderRepository interface {
	GetAll() ([]domain.Rider, error)
	Get(id string) (domain.Rider, error)
	Save(rider domain.Rider) (domain.Rider, error)
	Update(rider domain.Rider) (domain.Rider, error)
}
