package services

import (
	"context"
	"customer-service/config"
	"customer-service/internal/core/domain"
	"customer-service/internal/core/interfaces"
	"customer-service/pkg/rabbitmq"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/sdk/trace"
	"testing"
)

type RabbitMQPublisherTestSuite struct {
	suite.Suite
	TestRabbitMQ  *rabbitmq.RabbitMQ
	TestPublisher interfaces.MessageBusPublisher
	Cfg           *config.Config
	TestData      struct {
		Customer domain.Customer
		User     domain.User
	}
}

func (suite *RabbitMQPublisherTestSuite) SetupSuite() {
	cfgPath := "../../../test/customer.config"
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(errors.WithStack(err))
	}

	rmqServer, err := rabbitmq.NewRabbitMQ(cfg)

	if err != nil {
		panic(errors.WithStack(err))
	}

	tracer := trace.NewTracerProvider()

	rmqPublisher := NewRabbitMQPublisher(rmqServer, tracer, cfg)

	suite.Cfg = cfg
	suite.TestRabbitMQ = rmqServer
	suite.TestPublisher = rmqPublisher
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

func (suite *RabbitMQPublisherTestSuite) TestRabbitMQPublisher_CreateCustomer() {
	ch, err := suite.TestRabbitMQ.Connection.Channel()

	suite.NoError(err)

	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	suite.NoError(err)

	err = ch.QueueBind(
		queue.Name,
		"customer.create",
		suite.Cfg.RabbitMQ.Exchange,
		false,
		nil)
	if err != nil {
		return
	}

	suite.NoError(err)

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	suite.NoError(err)

	err = suite.TestPublisher.CreateCustomer(context.Background(), suite.TestData.Customer)

	suite.NoError(err)

	for msg := range msgs {
		suite.Equal("customer.create", msg.RoutingKey)

		var customer domain.Customer

		err = json.Unmarshal(msg.Body, &customer)
		suite.NoError(err)

		suite.Equal(suite.TestData.Customer, customer)

		err = msg.Ack(true)

		suite.NoError(err)

		err = ch.Close()

		suite.NoError(err)

		return
	}

}

func (suite *RabbitMQPublisherTestSuite) TestRabbitMQPublisher_UpdateServiceArea() {
	ch, err := suite.TestRabbitMQ.Connection.Channel()

	suite.NoError(err)

	queue, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)

	suite.NoError(err)

	err = ch.QueueBind(
		queue.Name,
		"customer.update.serviceArea",
		suite.Cfg.RabbitMQ.Exchange,
		false,
		nil)
	if err != nil {
		return
	}

	suite.NoError(err)

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	suite.NoError(err)

	err = suite.TestPublisher.UpdateServiceArea(context.Background(), suite.TestData.Customer)

	suite.NoError(err)

	for msg := range msgs {
		suite.Equal("customer.update.serviceArea", msg.RoutingKey)

		var message struct {
			ID          string
			ServiceArea int
		}

		err = json.Unmarshal(msg.Body, &message)
		suite.NoError(err)

		suite.Equal(suite.TestData.Customer.UserID, message.ID)
		suite.Equal(suite.TestData.Customer.ServiceArea, message.ServiceArea)

		err = msg.Ack(true)

		suite.NoError(err)

		err = ch.Close()

		suite.NoError(err)

		return
	}
}

func TestIntegration_RabbitMQPublisherTestSuite(t *testing.T) {
	testSuite := new(RabbitMQPublisherTestSuite)
	suite.Run(t, testSuite)
}
