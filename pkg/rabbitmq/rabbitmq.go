package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ(connstr string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(connstr)

	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	err = channel.ExchangeDeclare(
		"topics",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	err = channel.Qos(
		1,
		0,
		false,
	)

	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		Connection: conn,
		Channel:    channel,
	}, nil
}

func (r *RabbitMQ) Close() {
	r.Connection.Close()
}
