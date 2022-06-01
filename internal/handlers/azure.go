package handlers

import (
	"context"
	"customer-service/config"
	"customer-service/internal/core/domain"
	"customer-service/internal/core/interfaces"
	"customer-service/pkg/azure"
	"encoding/json"
	"fmt"
)

type azureHandler struct {
	serviceBus *azure.ServiceBus
	service    interfaces.CustomerService
	handlers   map[string]func(topic string, body []byte, handler *azureHandler) error
	config     *config.Config
	channel    chan bool
}

func NewAzure(serviceBus *azure.ServiceBus, service interfaces.CustomerService, config *config.Config) *azureHandler {
	return &azureHandler{
		serviceBus: serviceBus,
		service:    service,
		handlers: map[string]func(topic string, body []byte, handler *azureHandler) error{
			"user.create": userCreateOrUpdate,
			"user.update": userCreateOrUpdate,
		},
		config: config,
	}
}

func userCreateOrUpdate(topic string, body []byte, handler *azureHandler) error {
	var user domain.User

	if err := json.Unmarshal(body, &user); err != nil {
		return err
	}

	if err := handler.service.SaveOrUpdateUser(context.Background(), user); err != nil {
		return err
	}

	return nil
}

func (handler *azureHandler) Listen() {

	receiver, err := handler.serviceBus.Client.NewReceiverForQueue(
		"customer-service-queue",
		nil,
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	handler.channel = make(chan bool)

	go func() {
		for {
			select {
			case <-handler.channel:
				return
			default:
				msgs, err := receiver.ReceiveMessages(
					context.Background(),
					1,
					nil,
				)

				if err != nil {
					fmt.Println(err)
					return
				}

				for _, msg := range msgs {

					if msg.Subject != nil && *msg.Subject != "" {
						fun, exist := handler.handlers[*msg.Subject]

						if exist {
							err = fun(*msg.Subject, msg.Body, handler)
							if err == nil {
								_ = receiver.CompleteMessage(context.Background(), msg, nil)
								continue
							}
						}
					} else {
						fmt.Println("Message contains no subject: ", msg.MessageID)
						_ = receiver.CompleteMessage(context.Background(), msg, nil)
						continue
					}

					fmt.Println(err)
					_ = receiver.AbandonMessage(context.Background(), msg, nil)
				}
			}
		}
	}()
}

func (handler *azureHandler) Quit() {
	handler.channel <- true
}
