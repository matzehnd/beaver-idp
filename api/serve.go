package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupApi() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r

}

func Serve(port string) {
	r := setupApi()
	r.Run(":" + port)
}
