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

	dbConn := os.Getenv("DB")

	if dbConn == "" {
		panic("no db connection string")
	}

	privateKey := os.Getenv("PRIVATE_KEY")

	if privateKey == "" {
		panic("no private key")
	}

	cfg := config.LoadConfig(dbConn)
	config.InitDb(*cfg)
	defer cfg.DB.Close()
	eventStore := eventstore.NewPostgresEventStore(cfg.DB)

	tokenService := domain.NewTokenService([]byte(privateKey))
	userService := domain.NewUserService(eventStore)
	thingService := domain.NewThingService(eventStore)

	err := userService.RebuildEventStream() // Event Stream beim Start neu bilden
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	v1 := r.Group("/v1")

	http.NewV1Handler(v1, userService, tokenService, thingService)

	r.Run(":" + port)

}
