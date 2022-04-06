package customer_service

import (
	"customer-service/internal/core/domain"
	"customer-service/internal/core/ports"
	"errors"
)

type service struct {
	customerRepository ports.CustomerRepository
	messagePublisher   ports.MessageBusPublisher
}

func New(riderRepository ports.CustomerRepository, messagePublisher ports.MessageBusPublisher) *service {
	return &service{
		customerRepository: riderRepository,
		messagePublisher:   messagePublisher,
	}
}

func (srv *service) GetAll() ([]domain.Customer, error) {
	return srv.customerRepository.GetAll()
}

func (srv *service) Get(id string) (domain.Customer, error) {
	return srv.customerRepository.Get(id)
}

func (srv *service) Create(userId string, serviceArea int) (domain.Customer, error) {
	user, err := srv.customerRepository.GetUser(userId)

	customer := domain.NewCustomer(user, serviceArea)

	customer, err = srv.customerRepository.Save(customer)

	if err != nil {
		return domain.Customer{}, errors.New("saving new customer failed")
	}

	srv.messagePublisher.CreateCustomer(customer)
	return customer, nil
}

func (srv *service) UpdateServiceArea(id string, serviceArea int) (domain.Customer, error) {
	customer, err := srv.Get(id)

	if err != nil {
		return domain.Customer{}, errors.New("could not find customer with id")
	}

	customer.ServiceArea = serviceArea

	customer, err = srv.customerRepository.Update(customer)

	if err != nil {
		return domain.Customer{}, errors.New("saving new customer failed")
	}

	srv.messagePublisher.UpdateServiceArea(customer)
	return customer, nil
}

func (srv *service) SaveOrUpdateUser(user domain.User) error {
	if user.Name == "" || user.LastName == "" || user.ID == "" {
		return errors.New("incomplete user data")
	}

	err := srv.customerRepository.SaveOrUpdateCustomer(user)

	return err
}
