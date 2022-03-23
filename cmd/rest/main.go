package main

import (
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
	"os"
	"rider-service/docs"
	"rider-service/internal/core/services/ridersrv"
	"rider-service/internal/handlers/riderhdl"
	"rider-service/internal/repositories/riderrepo"

	"github.com/gin-gonic/gin"

	ginSwagger "github.com/swaggo/gin-swagger"
	_ "rider-service/docs"
)

const defaultPort = ":1234"

func main() {
	setupSwagger()

	riderRepository, err := riderrepo.NewCockroachDB("postgresql://root@localhost:26257/test?sslmode=disable")

	if err != nil {
		panic(err)
	}

	riderService := ridersrv.New(riderRepository)
	riderHandler := riderhdl.NewHTTPHandler(riderService)

	router := gin.New()

	router.GET("/riders", riderHandler.GetAll)
	router.GET("/riders/:id", riderHandler.Get)
	router.POST("/riders", riderHandler.Create)
	router.PUT("/riders/:id", riderHandler.UpdateRider)
	router.PUT("/riders/:id/location", riderHandler.UpdateLocation)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Fatal(router.Run(port))
}

func setupSwagger() {
	docs.SwaggerInfo.Title = "Rider service API"
	docs.SwaggerInfo.Description = "The rider service manages all riders for the BikePack system."
}
