package rabbitmq_service

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"rider-service/internal/core/domain"
	"rider-service/pkg/rabbitmq"
)

type rabbitmqPublisher rabbitmq.RabbitMQ

func NewRabbitMQPublisher(rabbitmq *rabbitmq.RabbitMQ) *rabbitmqPublisher {
	return &rabbitmqPublisher{Connection: rabbitmq.Connection, Channel: rabbitmq.Channel}
}

func (rmq *rabbitmqPublisher) CreateRider(rider domain.Rider) error {
	js, err := json.Marshal(rider)

	if err != nil {
		return err
	}

	err = rmq.Channel.Publish(
		"topics",
		"rider.create",
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

func (rmq *rabbitmqPublisher) UpdateRider(rider domain.Rider) error {
	js, err := json.Marshal(rider)

	if err != nil {
		return err
	}

	err = rmq.Channel.Publish(
		"topics",
		"rider.update",
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
