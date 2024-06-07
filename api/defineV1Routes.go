package api

import (
	"github.com/gin-gonic/gin"
)

func DefineV1Routes(rg *gin.RouterGroup) {

	rg.POST("/register", register)
}

type RegisterTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func register(c *gin.Context) {
	var data RegisterTO
	c.BindJSON(&data)

}
