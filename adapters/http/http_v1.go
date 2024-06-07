package http

import (
	"beaver/idp/core/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *domain.UserService
}

func NewV1Handler(router *gin.RouterGroup, userService *domain.UserService) {
	handler := &UserHandler{userService: userService}
	router.POST("/users", handler.createUser)
	router.GET("/users/:id", handler.getUser)
}

type RegisterUserTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) createUser(c *gin.Context) {
	var user RegisterUserTO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.userService.RegisterUser(domain.RegisterUser(user)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) getUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
