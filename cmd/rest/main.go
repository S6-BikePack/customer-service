package main

import (
	"customer-service/internal/core/services/customer_service"
	"customer-service/internal/core/services/rabbitmq_service"
	"customer-service/internal/handlers"
	"customer-service/internal/repositories"
	"customer-service/pkg/rabbitmq"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

const defaultPort = ":1234"
const defaultRmqConn = "amqp://user:password@localhost:5672/"
const defaultDbConn = "postgresql://user:password@localhost:5432/customer"

func main() {
	dbConn := GetEnvOrDefault("DATABASE", defaultDbConn)

	customerRepository, err := repositories.NewCockroachDB(dbConn)

	if err != nil {
		panic(err)
	}

	rmqConn := GetEnvOrDefault("RABBITMQ", defaultRmqConn)

	rmqServer, err := rabbitmq.NewRabbitMQ(rmqConn)

	if err != nil {
		panic(err)
	}

	rmqPublisher := rabbitmq_service.NewRabbitMQPublisher(rmqServer)

	customerService := customer_service.New(customerRepository, rmqPublisher)

	router := gin.New()

	customerHandler := handlers.NewRest(customerService, router)
	customerHandler.SetupEndpoints()
	customerHandler.SetupSwagger()

	port := GetEnvOrDefault("PORT", defaultPort)

	log.Fatal(router.Run(port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
