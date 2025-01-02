package main

import (
	"beaver/idp/adapters/eventstore"
	"beaver/idp/adapters/http"
	"beaver/idp/config"
	"beaver/idp/core/domain"
	"os"
	"time"

	"github.com/gin-contrib/cors"
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
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	r.Use(cors.New(config))
	v1 := r.Group("/v1")

	http.NewV1Handler(v1, userService, tokenService, thingService)

	r.Run(":" + port)

}
