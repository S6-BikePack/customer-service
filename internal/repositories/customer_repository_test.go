package repositories

import (
	"context"
	"customer-service/config"
	"customer-service/internal/core/domain"
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

type CustomerRepositoryTestSuite struct {
	suite.Suite
	TestDb   *gorm.DB
	TestRepo *customerRepository
	Cfg      *config.Config
	TestData struct {
		Customer domain.Customer
		User     domain.User
	}
}

func (suite *CustomerRepositoryTestSuite) SetupSuite() {
	cfgPath := "../../test/customer.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database)
	db, err := gorm.Open(postgres.Open(dsn))
	db.Debug()

	if err != nil {
		panic(errors.WithStack(err))
	}

	repository, err := NewCustomerRepository(db)

	if err != nil {
		panic(errors.WithStack(err))
	}

	db.Exec("DELETE FROM public.customers")
	db.Exec("DELETE FROM public.users")

	db.Exec("INSERT INTO public.users (id, name, last_name) VALUES ('test-id', 'test-name', 'test-lastname')")
	db.Exec("INSERT INTO public.customers (user_id, service_area) VALUES ('test-id', 1)")

	suite.Cfg = cfg
	suite.TestDb = db
	suite.TestRepo = repository
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

func (suite *CustomerRepositoryTestSuite) TestRepository_Get() {
	result, err := suite.TestRepo.Get(context.Background(), suite.TestData.Customer.UserID)

	suite.NoError(err)

	suite.EqualValues(suite.TestData.Customer, result)
}

func (suite *CustomerRepositoryTestSuite) TestRepository_Get_NotFound() {
	_, err := suite.TestRepo.Get(context.Background(), "test")

	suite.Error(err)
}

func (suite *CustomerRepositoryTestSuite) TestRepository_Save() {
	suite.TestDb.Exec("INSERT INTO public.users (id, name, last_name) VALUES ('test-id-2', 'test-name', 'test-lastname')")

	newCustomer := suite.TestData.Customer
	newCustomer.User = domain.User{}
	newCustomer.UserID = "test-id-2"

	_, err := suite.TestRepo.Save(context.Background(), newCustomer)

	suite.NoError(err)

	queryResult := domain.Customer{}
	suite.TestDb.Raw("SELECT * FROM public.customers WHERE user_id=?",
		newCustomer.UserID).Scan(&queryResult)

	suite.EqualValues(queryResult.ServiceArea, newCustomer.ServiceArea)
}

func (suite *CustomerRepositoryTestSuite) TestRepository_Update() {
	updated := suite.TestData.Customer
	updated.UserID = "test-id-2"
	updated.ServiceArea = 2
	updated.User = domain.User{}

	_, err := suite.TestRepo.Update(context.Background(), updated)

	suite.NoError(err)

	queryResult := domain.Customer{}
	suite.TestDb.Raw("SELECT * FROM public.customers WHERE user_id=?",
		updated.UserID).Scan(&queryResult)

	suite.EqualValues(queryResult.ServiceArea, updated.ServiceArea)
}

func (suite *CustomerRepositoryTestSuite) TestRepository_SaveUser() {
	user := domain.User{
		ID:       "test-id-3",
		Name:     "test-name-3",
		LastName: "test-lastname-3",
	}

	err := suite.TestRepo.SaveOrUpdateUser(context.Background(), user)

	suite.NoError(err)

	queryResult := domain.User{}
	suite.TestDb.Raw("SELECT * FROM public.users WHERE id=?",
		user.ID).Scan(&queryResult)

	suite.EqualValues(queryResult, user)
}

func (suite *CustomerRepositoryTestSuite) TestRepository_UpdateUser() {
	user := domain.User{
		ID:       "test-id-3",
		Name:     "new-name-3",
		LastName: "new-lastname-3",
	}

	err := suite.TestRepo.SaveOrUpdateUser(context.Background(), user)

	suite.NoError(err)

	queryResult := domain.User{}
	suite.TestDb.Raw("SELECT * FROM public.users WHERE id=?",
		user.ID).Scan(&queryResult)

	suite.EqualValues(queryResult, user)
}

func (suite *CustomerRepositoryTestSuite) TestRepository_GetUser() {
	user := domain.User{
		ID:       "test-id",
		Name:     "test-name",
		LastName: "test-lastname",
	}

	result, err := suite.TestRepo.GetUser(context.Background(), user.ID)

	suite.NoError(err)

	suite.EqualValues(result, user)
}

func TestIntegration_CustomerRepositoryTestSuite(t *testing.T) {
	repoSuite := new(CustomerRepositoryTestSuite)
	suite.Run(t, repoSuite)
}
