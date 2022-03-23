package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"rider-service/internal/core/services/ridersrv"
	"rider-service/internal/graph"
	"rider-service/internal/graph/generated"
	"rider-service/internal/repositories/riderrepo"
)

const defaultPort = ":1235"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := gin.Default()
	router.POST("/query", graphqlHandler())
	router.GET("/", playgroundHandler())
	log.Fatal(router.Run(port))
}

func graphqlHandler() gin.HandlerFunc {
	riderRepository, err := riderrepo.NewCockroachDB("postgresql://root@localhost:26257/test?sslmode=disable")

	if err != nil {
		panic(err)
	}

	riderService := ridersrv.New(riderRepository)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{RiderService: riderService}}))

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	srv := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}
