package customer_service

import (
	"customer-service/internal/core/domain"
	"customer-service/internal/core/ports"
	"errors"
	"github.com/google/uuid"
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

func (srv *service) Get(uuid uuid.UUID) (domain.Customer, error) {
	return srv.customerRepository.Get(uuid)
}

func (srv *service) Create(name, lastName, email string) (domain.Customer, error) {
	customer := domain.NewCustomer(name, lastName, email, 0)

	customer, err := srv.customerRepository.Save(customer)

	if err != nil {
		return domain.Customer{}, errors.New("saving new customer failed")
	}

	srv.messagePublisher.CreateCustomer(customer)
	return customer, nil
}

func (srv *service) UpdateCustomerDetails(uuid uuid.UUID, name, lastName, email string) (domain.Customer, error) {
	customer, err := srv.Get(uuid)

	if err != nil {
		return domain.Customer{}, errors.New("could not find customer with id")
	}

	if name != "" {
		customer.Name = name
	}

	if lastName != "" {
		customer.LastName = lastName
	}

	if email != "" {
		customer.Email = email
	}

	customer, err = srv.customerRepository.Update(customer)

	if err != nil {
		return domain.Customer{}, errors.New("saving new customer failed")
	}

	srv.messagePublisher.UpdateCustomerDetails(customer)
	return customer, nil
}

func (srv *service) UpdateServiceArea(uuid uuid.UUID, serviceArea int) (domain.Customer, error) {
	customer, err := srv.Get(uuid)

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
