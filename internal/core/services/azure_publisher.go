package services

import (
	"context"
	"customer-service/config"
	"customer-service/internal/core/domain"
	"customer-service/pkg/azure"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

type azurePublisher struct {
	serviceBus *azure.ServiceBus
	sender     *azservicebus.Sender
	config     *config.Config
}

func NewAzurePublisher(serviceBus *azure.ServiceBus, cfg *config.Config) *azurePublisher {
	return &azurePublisher{serviceBus: serviceBus, config: cfg}
}

func (az *azurePublisher) CreateCustomer(ctx context.Context, customer domain.Customer) error {
	return az.publishJson(ctx, "create", customer)
}

func (az *azurePublisher) UpdateServiceArea(ctx context.Context, customer domain.Customer) error {
	var body = struct {
		ID          string
		ServiceArea int
	}{
		ID:          customer.UserID,
		ServiceArea: customer.ServiceArea,
	}

	return az.publishJson(ctx, "update.serviceArea", body)
}

func (az *azurePublisher) publishJson(ctx context.Context, topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	topic = fmt.Sprintf("customer.%s", topic)

	sender, err := az.serviceBus.Client.NewSender(topic, nil)

	defer func(sender *azservicebus.Sender, ctx context.Context) {
		_ = sender.Close(ctx)
	}(sender, ctx)

	if err != nil {
		return err
	}

	err = sender.SendMessage(ctx, &azservicebus.Message{
		Body:    js,
		Subject: &topic,
	}, nil)

	if err != nil {
		return err
	}

	return err
}
