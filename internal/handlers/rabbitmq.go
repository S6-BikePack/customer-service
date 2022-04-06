package handlers

import (
	"customer-service/internal/core/domain"
	"customer-service/internal/core/ports"
	"customer-service/pkg/rabbitmq"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
)

type rabbitmqHandler struct {
	rabbitmq *rabbitmq.RabbitMQ
	service  ports.CustomerService
	handlers map[string]func(topic string, body []byte, handler *rabbitmqHandler) error
}

func NewRabbitMQ(rabbitmq *rabbitmq.RabbitMQ, service ports.CustomerService) *rabbitmqHandler {
	return &rabbitmqHandler{
		rabbitmq: rabbitmq,
		service:  service,
		handlers: map[string]func(topic string, body []byte, handler *rabbitmqHandler) error{
			"user.create": UserCreateOrUpdate,
			"user.update": UserCreateOrUpdate,
		},
	}
}

func UserCreateOrUpdate(topic string, body []byte, handler *rabbitmqHandler) error {
	var user domain.User

	if err := json.Unmarshal(body, &user); err != nil {
		return err
	}

	if err := handler.service.SaveOrUpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (handler *rabbitmqHandler) Listen() {

	q, err := handler.rabbitmq.Channel.QueueDeclare(
		"customerQueue",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	for _, s := range maps.Keys(handler.handlers) {
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
			fun, exist := handler.handlers[msg.RoutingKey]

			if exist {
				err = fun(msg.RoutingKey, msg.Body, handler)
				if err == nil {
					msg.Ack(false)
					continue
				}
			}

			fmt.Println(err)
			msg.Nack(false, true)
		}
	}()

	<-forever
}

type MessageHandler struct {
	topic   string
	handler func(topic string, body []byte, handler *rabbitmqHandler) error
}
