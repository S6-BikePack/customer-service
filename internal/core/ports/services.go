package ports

import "rider-service/internal/core/domain"

type RiderService interface {
	GetAll() ([]domain.Rider, error)
	Get(id string) (domain.Rider, error)
	Create(name string, status string) (domain.Rider, error)
	Update(id string, name string, status string) (domain.Rider, error)
	UpdateLocation(id string, location domain.Location) (domain.Rider, error)
}
