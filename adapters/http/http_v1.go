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
	router.POST("/users/token", handler.createToken)
	router.GET("/users/:id", handler.getUser)
	router.GET("/validation", handler.validate)
}

type RegisterUserTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type CreateTokenTO struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ValidityInHours *int   `json:"validityInHours,omitempty"`
}

type TokenTO struct {
	Token string `json:"token"`
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

func (h *UserHandler) createToken(c *gin.Context) {
	var user CreateTokenTO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.userService.CreateToken(domain.CreateToken(user))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) validate(c *gin.Context) {
	var token TokenTO
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	validated, err := h.userService.ValidateToken(token.Token)

	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	if validated.Valid {
		c.Status(http.StatusOK)
		return
	}
	c.Status(http.StatusUnauthorized)
}
