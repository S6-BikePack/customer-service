package services

import (
	"context"
	"customer-service/config"
	"customer-service/internal/core/domain"
	"customer-service/pkg/rabbitmq"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type rabbitmqPublisher struct {
	rabbitmq *rabbitmq.RabbitMQ
	tracer   trace.Tracer
	config   *config.Config
}

func NewRabbitMQPublisher(rabbitmq *rabbitmq.RabbitMQ, tracerProvider trace.TracerProvider, cfg *config.Config) *rabbitmqPublisher {
	return &rabbitmqPublisher{rabbitmq: rabbitmq, tracer: tracerProvider.Tracer("RabbitMQ.Publisher"), config: cfg}
}

func (rmq *rabbitmqPublisher) CreateCustomer(ctx context.Context, customer domain.Customer) error {
	var body = struct {
		ID          string
		ServiceArea int
	}{
		ID:          customer.UserID,
		ServiceArea: customer.ServiceArea,
	}

	return rmq.publishJson(ctx, "create", body)
}

func (rmq *rabbitmqPublisher) UpdateServiceArea(ctx context.Context, customer domain.Customer) error {
	return rmq.publishJson(ctx, "update.serviceArea", customer)
}

func (rmq *rabbitmqPublisher) publishJson(ctx context.Context, topic string, body interface{}) error {
	js, err := json.Marshal(body)

	if err != nil {
		return err
	}

	_, span := rmq.tracer.Start(ctx, "publish")

	span.AddEvent(
		"Published message to rabbitmq",
		trace.WithAttributes(
			attribute.String("topic", topic),
			attribute.String("body", string(js))))
	span.End()

	err = rmq.rabbitmq.Channel.Publish(
		rmq.config.RabbitMQ.Exchange,
		fmt.Sprintf("customer.%s", topic),
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
