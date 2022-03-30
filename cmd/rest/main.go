package main

import (
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
	"os"
	"rider-service/docs"
	"rider-service/internal/core/services/rabbitmq_service"
	"rider-service/internal/core/services/ridersrv"
	"rider-service/internal/handlers"
	"rider-service/internal/handlers/riderhdl"
	"rider-service/internal/repositories/riderrepo"
	"rider-service/pkg/rabbitmq"

	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"
	_ "rider-service/docs"
)

const defaultPort = ":1234"
const defaultRmqConn = "amqp://user:password@localhost:5672/"
const defaultDbConn = "postgresql://user:password@localhost:5432/rider"

func main() {
	setupSwagger()

	dbConn := os.Getenv("DATABASE")
	if dbConn == "" {
		dbConn = defaultDbConn
	}

	riderRepository, err := riderrepo.NewCockroachDB(dbConn)

	if err != nil {
		panic(err)
	}

	rmqConn := os.Getenv("RABBITMQ")
	if rmqConn == "" {
		rmqConn = defaultRmqConn
	}

	rmqServer, err := rabbitmq.NewRabbitMQ(rmqConn)

	if err != nil {
		panic(err)
	}

	rmqPublisher := rabbitmq_service.NewRabbitMQPublisher(rmqServer)

	riderService := ridersrv.New(riderRepository, rmqPublisher)

	rmqSubscriber := handlers.NewRabbitMQ(rmqServer, riderService)

	riderHandler := riderhdl.NewHTTPHandler(riderService)

	router := gin.New()

	api := router.Group("/api")
	api.GET("/riders", riderHandler.GetAll)
	api.GET("/riders/:id", riderHandler.Get)
	api.POST("/riders", riderHandler.Create)
	api.PUT("/riders/:id", riderHandler.UpdateRider)
	api.PUT("/riders/:id/location", riderHandler.UpdateLocation)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	go rmqSubscriber.Listen("delivery.#")
	log.Fatal(router.Run(port))
}

func setupSwagger() {
	docs.SwaggerInfo.Title = "Rider service API"
	docs.SwaggerInfo.Description = "The rider service manages all riders for the BikePack system."
}
