package handlers

import (
	"context"
	"customer-service/config"
	"customer-service/internal/core/domain"
	"customer-service/internal/core/interfaces"
	"customer-service/pkg/rabbitmq"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/maps"
)

type rabbitmqHandler struct {
	rabbitmq *rabbitmq.RabbitMQ
	service  interfaces.CustomerService
	handlers map[string]func(topic string, body []byte, handler *rabbitmqHandler) error
	config   *config.Config
	channel  chan bool
}

func NewRabbitMQ(rabbitmq *rabbitmq.RabbitMQ, service interfaces.CustomerService, config *config.Config) *rabbitmqHandler {
	return &rabbitmqHandler{
		rabbitmq: rabbitmq,
		service:  service,
		handlers: map[string]func(topic string, body []byte, handler *rabbitmqHandler) error{
			"user.create": UserCreateOrUpdate,
			"user.update": UserCreateOrUpdate,
		},
		config: config,
	}
}

func UserCreateOrUpdate(topic string, body []byte, handler *rabbitmqHandler) error {
	var user domain.User

	if err := json.Unmarshal(body, &user); err != nil {
		return err
	}

	if err := handler.service.SaveOrUpdateUser(context.Background(), user); err != nil {
		return err
	}

	return nil
}

func (handler *rabbitmqHandler) Listen() {

	q, err := handler.rabbitmq.Channel.QueueDeclare(
		handler.config.Server.Service+"Queue",
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
			handler.config.RabbitMQ.Exchange,
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

	handler.channel = make(chan bool)

	go func() {
		for msg := range msgs {
			select {
			case <-handler.channel:
				return
			default:
				fun, exist := handler.handlers[msg.RoutingKey]

				if exist {
					err = fun(msg.RoutingKey, msg.Body, handler)
					if err == nil {
						_ = msg.Ack(false)
						continue
					}
				}

				fmt.Println(err)
				_ = msg.Nack(false, true)
			}

		}
	}()
}

func (handler *rabbitmqHandler) Quit() {
	handler.channel <- true
}

type MessageHandler struct {
	topic   string
	handler func(topic string, body []byte, handler *rabbitmqHandler) error
}
