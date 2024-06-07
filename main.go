package main

import (
	"beaver/idp/adapters/eventstore"
	"beaver/idp/adapters/http"
	"beaver/idp/config"
	"beaver/idp/core/domain"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		panic("no port defined")
	}

	cfg := config.LoadConfig()
	eventStore := eventstore.NewPostgresEventStore(cfg.DB)

	userService := domain.NewUserService(eventStore)
	userService.RebuildEventStream() // Event Stream beim Start neu bilden

	r := gin.Default()
	v1 := r.Group("/v1")

	http.NewV1Handler(v1, userService)

	r.Run(":" + port)

}
