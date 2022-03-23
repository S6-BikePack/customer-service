package ridersrv

import (
	"errors"
	"rider-service/internal/core/domain"
	"rider-service/internal/core/ports"
)

type service struct {
	riderRepository ports.RiderRepository
}

func New(riderRepository ports.RiderRepository) *service {
	return &service{
		riderRepository: riderRepository,
	}
}

func (srv *service) GetAll() ([]domain.Rider, error) {
	return srv.riderRepository.GetAll()
}

func (srv *service) Get(id string) (domain.Rider, error) {
	return srv.riderRepository.Get(id)
}

func (srv *service) Create(name string, status string) (domain.Rider, error) {
	rider := domain.NewRider(name, status, domain.Location{})

	rider, err := srv.riderRepository.Save(rider)

	if err != nil {
		return domain.Rider{}, errors.New("saving new rider failed")
	}

	return rider, nil
}

func (srv *service) Update(id string, name string, status string) (domain.Rider, error) {
	rider, err := srv.Get(id)

	if err != nil {
		return domain.Rider{}, errors.New("could not find rider with id")
	}

	if status != "" {
		rider.Status = status
	}

	if name != "" {
		rider.Name = name
	}

	rider, err = srv.riderRepository.Update(rider)

	if err != nil {
		return domain.Rider{}, errors.New("saving new rider failed")
	}

	return rider, nil
}

func (srv *service) UpdateLocation(id string, location domain.Location) (domain.Rider, error) {
	rider, err := srv.Get(id)

	if err != nil {
		return domain.Rider{}, errors.New("could not find rider with id")
	}

	rider.Location = location

	rider, err = srv.riderRepository.Update(rider)

	if err != nil {
		return domain.Rider{}, errors.New("saving new rider failed")
	}

	return rider, nil
}
