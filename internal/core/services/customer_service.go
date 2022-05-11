package services

import (
	"context"
	"customer-service/internal/core/domain"
	"customer-service/internal/core/interfaces"
	"errors"
)

type service struct {
	customerRepository interfaces.CustomerRepository
	messagePublisher   interfaces.MessageBusPublisher
}

func NewCustomerService(riderRepository interfaces.CustomerRepository, messagePublisher interfaces.MessageBusPublisher) *service {
	return &service{
		customerRepository: riderRepository,
		messagePublisher:   messagePublisher,
	}
}

func (srv *service) GetAll(ctx context.Context) ([]domain.Customer, error) {
	return srv.customerRepository.GetAll(ctx)
}

func (srv *service) Get(ctx context.Context, id string) (domain.Customer, error) {
	return srv.customerRepository.Get(ctx, id)
}

func (srv *service) Create(ctx context.Context, userId string, serviceArea int) (domain.Customer, error) {
	user, err := srv.customerRepository.GetUser(ctx, userId)

	if err != nil {
		return domain.Customer{}, err
	}

	customer := domain.NewCustomer(user, serviceArea)

	customer, err = srv.customerRepository.Save(ctx, customer)

	if err != nil {
		return domain.Customer{}, errors.New("saving new customer failed")
	}

	_ = srv.messagePublisher.CreateCustomer(ctx, customer)
	return customer, nil
}

func (srv *service) UpdateServiceArea(ctx context.Context, id string, serviceArea int) (domain.Customer, error) {
	customer, err := srv.Get(ctx, id)
	updated := customer

	if err != nil {
		return domain.Customer{}, errors.New("could not find customer with id")
	}

	updated.ServiceArea = serviceArea

	updated, err = srv.customerRepository.Update(ctx, updated)

	if err != nil {
		return customer, errors.New("saving new customer failed")
	}

	_ = srv.messagePublisher.UpdateServiceArea(ctx, updated)
	return updated, nil
}

func (srv *service) SaveOrUpdateUser(ctx context.Context, user domain.User) error {
	if user.Name == "" || user.LastName == "" || user.ID == "" {
		return errors.New("incomplete user data")
	}

	err := srv.customerRepository.SaveOrUpdateUser(ctx, user)

	return err
}
