package rabbitmq_service

import (
	"customer-service/internal/core/domain"
	"customer-service/pkg/rabbitmq"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitmqPublisher rabbitmq.RabbitMQ

func NewRabbitMQPublisher(rabbitmq *rabbitmq.RabbitMQ) *rabbitmqPublisher {
	return &rabbitmqPublisher{Connection: rabbitmq.Connection, Channel: rabbitmq.Channel}
}

func (rmq *rabbitmqPublisher) CreateCustomer(customer domain.Customer) error {
	js, err := json.Marshal(customer)

	if err != nil {
		return err
	}

	err = rmq.publishMessage("customer.create", js)

	return err
}

func (rmq *rabbitmqPublisher) UpdateServiceArea(customer domain.Customer) error {
	js, err := json.Marshal(customer)

	if err != nil {
		return err
	}

	err = rmq.publishMessage("customer.update.serviceArea", js)

	return err
}

func (rmq *rabbitmqPublisher) publishMessage(key string, body []byte) error {
	err := rmq.Channel.Publish(
		"topics",
		key,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		},
	)

	return err
}
