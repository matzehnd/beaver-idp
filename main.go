package main

import (
	"beaver/idp/adapters/eventstore"
	"beaver/idp/adapters/http"
	"beaver/idp/config"
	"beaver/idp/core/domain"
	"fmt"
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

	fmt.Println(privateKey)

	cfg := config.LoadConfig(dbConn)
	defer cfg.DB.Close()
	eventStore := eventstore.NewPostgresEventStore(cfg.DB)

	userService := domain.NewUserService(eventStore, []byte(privateKey))
	userService.RebuildEventStream() // Event Stream beim Start neu bilden

	r := gin.Default()
	v1 := r.Group("/v1")

	http.NewV1Handler(v1, userService)

	r.Run(":" + port)

}
