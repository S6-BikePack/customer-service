package main

import (
	"context"
	"customer-service/config"
	"customer-service/internal/core/services"
	"customer-service/internal/handlers"
	"customer-service/internal/repositories"
	"customer-service/pkg/azure"
	"customer-service/pkg/logging"
	"customer-service/pkg/tracing"
	"fmt"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"

	"github.com/gin-gonic/gin"
)

const defaultConfig = "./config/local.config"

func main() {
	cfgPath := GetEnvOrDefault("config", defaultConfig)
	cfg, err := config.UseConfig(cfgPath)

	if err != nil {
		panic(err)
	}

	//--------------------------------------------------------------------------------------
	// Setup Logging and Tracing
	//--------------------------------------------------------------------------------------

	logger, err := logging.NewSimpleLogger(cfg)

	if err != nil {
		panic(err)
	}

	tracer, err := tracing.NewOpenTracing(cfg.Server.Service, cfg.Tracing.Host, cfg.Tracing.Port)

	if err != nil {
		panic(err)
	}

	//--------------------------------------------------------------------------------------
	// Setup Database
	//--------------------------------------------------------------------------------------

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Database, cfg.Database.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	if cfg.Database.Debug {
		db.Debug()
	}

	if err = db.Use(otelgorm.NewPlugin(otelgorm.WithTracerProvider(tracer))); err != nil {
		panic(err)
	}

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	customerRepository, err := repositories.NewCustomerRepository(db)

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	//--------------------------------------------------------------------------------------
	// Setup Azure Service Bus
	//--------------------------------------------------------------------------------------

	azServiceBus, err := azure.NewAzureServiceBus(cfg)

	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	azPublisher := services.NewAzurePublisher(azServiceBus, cfg)

	//--------------------------------------------------------------------------------------
	// Setup Services
	//--------------------------------------------------------------------------------------

	customerService := services.NewCustomerService(customerRepository, azPublisher)

	//--------------------------------------------------------------------------------------
	// Setup HTTP server
	//--------------------------------------------------------------------------------------

	router := gin.New()
	router.Use(otelgin.Middleware(cfg.Server.Service, otelgin.WithTracerProvider(tracer)))

	deliveryHandler := handlers.NewRest(customerService, router, logger, cfg)
	deliveryHandler.SetupEndpoints()
	deliveryHandler.SetupSwagger()
	deliveryHandler.SetupHealthprobe()

	azSubscriber := handlers.NewAzure(azServiceBus, customerService, cfg)

	go azSubscriber.Listen()
	logger.Fatal(context.Background(), router.Run(cfg.Server.Port))
}

func GetEnvOrDefault(environmentKey, defaultValue string) string {
	returnValue := os.Getenv(environmentKey)
	if returnValue == "" {
		returnValue = defaultValue
	}
	return returnValue
}
