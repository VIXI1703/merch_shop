package handlers

import (
	"github.com/gin-gonic/gin"
	"merch_shop/internal/model"
	"net/http"
)

type userService interface {
	Authenticate(username, password string) (string, error)
}

type AuthHandler struct {
	userService userService
}

func NewAuthHandler(userService userService) AuthHandler {
	return AuthHandler{
		userService: userService,
	}
}

func (handler *AuthHandler) Routes(c *gin.RouterGroup) {
	c.POST("/auth", handler.Authenticate)
}

func (h AuthHandler) Authenticate(c *gin.Context) {
	var request model.AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{"Required fields are empty or not valid"})
		return
	}
	token, err := h.userService.Authenticate(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.AuthResponse{Token: token})
}
