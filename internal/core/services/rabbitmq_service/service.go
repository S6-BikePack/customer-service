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
	var body = struct {
		ID          string
		ServiceArea int
	}{
		ID:          customer.UserID,
		ServiceArea: customer.ServiceArea,
	}

	return rmq.publishJson("customer.create", body)

}

func (rmq *rabbitmqPublisher) UpdateServiceArea(customer domain.Customer) error {
	return rmq.publishJson("customer.update.serviceArea", customer)
}

func (rmq *rabbitmqPublisher) publishJson(topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	err = rmq.Channel.Publish(
		"topics",
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
