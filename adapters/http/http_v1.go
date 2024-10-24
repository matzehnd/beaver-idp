package http

import (
	"beaver/idp/core/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService  *domain.UserService
	tokenService *domain.TokenService
	thingService *domain.ThingService
}

func NewV1Handler(router *gin.RouterGroup, userService *domain.UserService, tokenService *domain.TokenService, thingService *domain.ThingService) {
	handler := &UserHandler{userService: userService, tokenService: tokenService, thingService: thingService}
	router.POST("/users", handler.registerUser)
	router.POST("/users/token", handler.createUserToken)
	router.GET("/users/:id", handler.getUser)
	router.GET("/validation", handler.validate)
	router.POST("things", handler.registerThing)
	router.POST("things/token", handler.createThingToken)
}

type RegisterUserTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterThingTO struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}
type CreateUserTokenTO struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ValidityInHours *int   `json:"validityInHours,omitempty"`
}
type CreateThingTokenTO struct {
	Id              string `json:"id"`
	Password        string `json:"password"`
	ValidityInHours *int   `json:"validityInHours,omitempty"`
}

type TokenTO struct {
	Token string `json:"token"`
}

func (h *UserHandler) registerUser(c *gin.Context) {
	var user RegisterUserTO
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.userService.RegisterUser(domain.RegisterUser(user)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"email": user.Email})
}

func (h *UserHandler) registerThing(c *gin.Context) {
	var thing RegisterThingTO
	if err := c.ShouldBindJSON(&thing); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.thingService.RegisterThing(domain.RegisterThing(thing)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": thing.Id})
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

func (h *UserHandler) createUserToken(c *gin.Context) {
	var query CreateUserTokenTO
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userService.GetUser(query.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	if isValid := h.userService.PasswordIsValid(*user, query.Password); !isValid {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	token, err := h.tokenService.CreateToken(user.Email, query.ValidityInHours, h.userService.UserIsAdmin(*user))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) createThingToken(c *gin.Context) {
	var query CreateThingTokenTO
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	thing, err := h.thingService.GetThing(query.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	if isValid := h.thingService.PasswordIsValid(*thing, query.Password); !isValid {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	token, err := h.tokenService.CreateToken(thing.Id, query.ValidityInHours, false)
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
	validated, err := h.tokenService.ValidateToken(token.Token)

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
