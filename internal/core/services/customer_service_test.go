package services

import (
	"context"
	"customer-service/internal/core/domain"
	"customer-service/internal/core/interfaces"
	"customer-service/internal/mock"
	"errors"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CustomerServiceTestSuite struct {
	suite.Suite
	MockRepository *mock.CustomerRepository
	MockPublisher  *mock.MessageBusPublisher
	TestService    interfaces.CustomerService
	TestData       struct {
		Customer domain.Customer
		User     domain.User
	}
}

func (suite *CustomerServiceTestSuite) SetupSuite() {
	repository := new(mock.CustomerRepository)
	publisher := new(mock.MessageBusPublisher)

	srv := NewCustomerService(repository, publisher)

	suite.MockRepository = repository
	suite.MockPublisher = publisher
	suite.TestService = srv
	suite.TestData = struct {
		Customer domain.Customer
		User     domain.User
	}{
		Customer: domain.Customer{
			UserID: "test-id",
			User: domain.User{
				ID:       "test-id",
				Name:     "test-name",
				LastName: "test-lastname",
			},
			ServiceArea: 1,
		},
		User: domain.User{
			ID:       "test-id-2",
			Name:     "test-name-2",
			LastName: "test-lastname-2",
		},
	}
}

func (suite *CustomerServiceTestSuite) SetupTest() {
	suite.MockPublisher.ExpectedCalls = nil
	suite.MockRepository.ExpectedCalls = nil
}

func (suite *CustomerServiceTestSuite) TestCustomerService_GetAll() {
	suite.MockRepository.On("GetAll").Return([]domain.Customer{suite.TestData.Customer}, nil)

	result, err := suite.TestService.GetAll(context.Background())

	suite.NoError(err)

	suite.MockRepository.AssertCalled(suite.T(), "GetAll")
	suite.Equal(1, len(result))
	suite.EqualValues(suite.TestData.Customer, result[0])
}

func (suite *CustomerServiceTestSuite) TestCustomerService_Get() {
	suite.MockRepository.On("Get", suite.TestData.Customer.UserID).Return(suite.TestData.Customer, nil)

	result, err := suite.TestService.Get(context.Background(), suite.TestData.Customer.UserID)

	suite.NoError(err)

	suite.MockRepository.AssertCalled(suite.T(), "Get", suite.TestData.Customer.UserID)
	suite.EqualValues(suite.TestData.Customer, result)
}

func (suite *CustomerServiceTestSuite) TestCustomerService_Get_NotFound() {
	suite.MockRepository.On("Get", suite.TestData.Customer.UserID).Return(domain.Customer{}, errors.New("could not find customer"))

	result, err := suite.TestService.Get(context.Background(), suite.TestData.Customer.UserID)

	suite.Error(err)
	suite.EqualValues(domain.Customer{}, result)

	suite.MockRepository.AssertCalled(suite.T(), "Get", suite.TestData.Customer.UserID)
}

func (suite *CustomerServiceTestSuite) TestCustomerService_Create() {
	suite.MockRepository.On("GetUser", suite.TestData.Customer.UserID).Return(suite.TestData.Customer.User, nil)
	suite.MockRepository.On("Save", mock2.Anything).Return(suite.TestData.Customer, nil)
	suite.MockPublisher.On("CreateCustomer", suite.TestData.Customer).Return(nil)

	result, err := suite.TestService.Create(context.Background(), suite.TestData.Customer.UserID, suite.TestData.Customer.ServiceArea)

	suite.NoError(err)

	suite.MockPublisher.AssertCalled(suite.T(), "CreateCustomer", suite.TestData.Customer)
	suite.EqualValues(suite.TestData.Customer, result)
}

func (suite *CustomerServiceTestSuite) TestCustomerService_Create_UserNotFound() {
	suite.MockRepository.On("GetUser", suite.TestData.Customer.UserID).Return(domain.User{}, errors.New("user not found"))

	_, err := suite.TestService.Create(context.Background(), suite.TestData.Customer.UserID, suite.TestData.Customer.ServiceArea)

	suite.MockRepository.AssertNotCalled(suite.T(), "Save")
	suite.Error(err)
}

func (suite *CustomerServiceTestSuite) TestCustomerService_Create_CouldNotSave() {
	suite.MockRepository.On("GetUser", suite.TestData.Customer.UserID).Return(suite.TestData.Customer.User, nil)
	suite.MockRepository.On("Save", mock2.Anything).Return(domain.Customer{}, errors.New("could not save customer"))
	suite.MockPublisher.On("CreateCustomer", suite.TestData.Customer).Return(nil)

	_, err := suite.TestService.Create(context.Background(), suite.TestData.Customer.UserID, suite.TestData.Customer.ServiceArea)

	suite.Error(err)

	suite.MockPublisher.AssertNotCalled(suite.T(), "CreateCustomer")
}

func (suite *CustomerServiceTestSuite) TestCustomerService_UpdateServiceArea() {
	updated := suite.TestData.Customer
	updated.ServiceArea = 2

	suite.MockRepository.On("Get", suite.TestData.Customer.UserID).Return(suite.TestData.Customer, nil)
	suite.MockRepository.On("Update", updated).Return(updated, nil)
	suite.MockPublisher.On("UpdateServiceArea", updated).Return(nil)

	result, err := suite.TestService.UpdateServiceArea(context.Background(), suite.TestData.Customer.UserID, updated.ServiceArea)

	suite.NoError(err)

	suite.EqualValues(updated, result)
}

func (suite *CustomerServiceTestSuite) TestCustomerService_UpdateServiceArea_CustomerNotFound() {
	updated := suite.TestData.Customer
	updated.ServiceArea = 2

	suite.MockRepository.On("Get", suite.TestData.Customer.UserID).Return(domain.Customer{}, errors.New("customer not found"))

	_, err := suite.TestService.UpdateServiceArea(context.Background(), suite.TestData.Customer.UserID, updated.ServiceArea)

	suite.Error(err)
}

func (suite *CustomerServiceTestSuite) TestCustomerService_UpdateServiceArea_CouldNotUpdate() {
	updated := suite.TestData.Customer
	updated.ServiceArea = 2

	suite.MockRepository.On("Get", suite.TestData.Customer.UserID).Return(suite.TestData.Customer, nil)
	suite.MockRepository.On("Update", updated).Return(suite.TestData.Customer, errors.New("could not update customer"))

	result, err := suite.TestService.UpdateServiceArea(context.Background(), suite.TestData.Customer.UserID, updated.ServiceArea)

	suite.Error(err)
	suite.EqualValues(suite.TestData.Customer, result)
}

func TestUnit_CustomerServiceTestSuite(t *testing.T) {
	repoSuite := new(CustomerServiceTestSuite)
	suite.Run(t, repoSuite)
}
