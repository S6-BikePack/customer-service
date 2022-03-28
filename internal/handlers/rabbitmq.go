package handlers

import (
	"errors"
	"fmt"
	"rider-service/internal/core/ports"
	"rider-service/pkg/rabbitmq"
)

type rabbitmqHandler struct {
	rabbitmq *rabbitmq.RabbitMQ
	service  ports.RiderService
}

func NewRabbitMQ(rabbitmq *rabbitmq.RabbitMQ, service ports.RiderService) *rabbitmqHandler {
	return &rabbitmqHandler{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (handler *rabbitmqHandler) Listen(topics ...string) {

	q, err := handler.rabbitmq.Channel.QueueDeclare(
		"riderQueue",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	for _, s := range topics {
		err = handler.rabbitmq.Channel.QueueBind(
			q.Name,
			s,
			"topics",
			false,
			nil)
		if err != nil {
			return
		}
	}

	msgs, err := handler.rabbitmq.Channel.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			err = handler.handleMessage(msg.RoutingKey, msg.Body)

			if err != nil {
				msg.Ack(false)
				continue
			}

			msg.Ack(true)
		}
	}()

	<-forever
}

func (handler *rabbitmqHandler) handleMessage(routing string, body []byte) error {
	switch routing {
	case "rider.create":
		fmt.Printf("Received message: %s\n", routing)
		fmt.Printf("Message payload: %s\n", body)
		return nil
	case "rider.update":
		return nil
	case "delivery.create":
		fmt.Printf("Route create")
		fmt.Printf("Message payload: %s\n", body)
		return nil
	case "delivery.update":
		fmt.Printf("Route update")
		fmt.Printf("Message payload: %s\n", body)
		return nil
	default:
		return errors.New("could not handle message")
	}
}
