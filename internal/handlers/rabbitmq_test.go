package handlers

import (
	"customer-service/config"
	"customer-service/internal/core/domain"
	"customer-service/internal/mock"
	"customer-service/pkg/rabbitmq"
	"encoding/json"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type RabbitMQHandlerTestSuite struct {
	suite.Suite
	MockCustomerService *mock.CustomerService
	TestRabbitMQ        *rabbitmq.RabbitMQ
	TestHandler         *rabbitmqHandler
	Cfg                 *config.Config
	TestData            struct {
		User     domain.User
		Customer domain.Customer
	}
}

func (suite *RabbitMQHandlerTestSuite) SetupSuite() {
	cfgPath := "../../test/customer.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	mockCustomerService := new(mock.CustomerService)

	rabbitMQ, err := rabbitmq.NewRabbitMQ(cfg)

	if err != nil {
		panic(errors.WithStack(err))
	}

	handler := NewRabbitMQ(rabbitMQ, mockCustomerService, cfg)

	go handler.Listen()

	time.Sleep(5 * time.Second) // wait for rabbitmq to start

	suite.Cfg = cfg
	suite.MockCustomerService = mockCustomerService
	suite.TestHandler = handler
	suite.TestRabbitMQ = rabbitMQ
	suite.TestData = struct {
		User     domain.User
		Customer domain.Customer
	}{
		User: domain.User{
			ID:       "test-id-2",
			Name:     "test-name-2",
			LastName: "test-lastname-2",
		},
		Customer: domain.Customer{
			UserID: "test-id",
			User: domain.User{
				ID:       "test-id",
				Name:     "test-name",
				LastName: "test-lastname",
			},
			ServiceArea: 1,
		},
	}
}

func (suite *RabbitMQHandlerTestSuite) SetupTest() {
	suite.MockCustomerService.ExpectedCalls = nil
}

func (suite *RabbitMQHandlerTestSuite) TearDownSuite() {
	suite.TestRabbitMQ.Close()
}

func (suite *RabbitMQHandlerTestSuite) TestHandler_UserCreateOrUpdate_Create() {
	suite.MockCustomerService.On("SaveOrUpdateUser", suite.TestData.User).Return(nil)

	err := publishJson(suite.TestRabbitMQ, suite.Cfg.RabbitMQ.Exchange, "user.create", suite.TestData.User)

	suite.NoError(err)

	for len(suite.MockCustomerService.Calls) < 1 {
	}

	suite.MockCustomerService.AssertCalled(suite.T(), "SaveOrUpdateUser", suite.TestData.User)
}

func (suite *RabbitMQHandlerTestSuite) TestHandler_UserCreateOrUpdate_Update() {
	suite.MockCustomerService.On("SaveOrUpdateUser", suite.TestData.User).Return(nil)

	err := publishJson(suite.TestRabbitMQ, suite.Cfg.RabbitMQ.Exchange, "user.update", suite.TestData.User)

	suite.NoError(err)

	for len(suite.MockCustomerService.Calls) < 1 {
	}

	suite.MockCustomerService.AssertCalled(suite.T(), "SaveOrUpdateUser", suite.TestData.User)
}

func publishJson(rabbitmq *rabbitmq.RabbitMQ, exchange, topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	err = rabbitmq.Channel.Publish(
		exchange,
		topic,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         js,
		},
	)

	return err
}

func TestIntegration_RabbitMQHandlerTestSuite(t *testing.T) {
	repoSuite := new(RabbitMQHandlerTestSuite)
	suite.Run(t, repoSuite)
}
