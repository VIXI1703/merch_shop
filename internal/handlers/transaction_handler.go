package handlers

import (
	"github.com/gin-gonic/gin"
	"merch_shop/internal/middleware"
	"merch_shop/internal/model"
	"net/http"
)

type transactionService interface {
	GetInfo(userId uint) (model.InfoResponse, error)
	SendCoin(userId uint, toUser string, amount uint) error
	BuyItem(userId uint, name string) error
}

type TransactionHandler struct {
	transactionService transactionService
}

func NewTransactionHandler(infoService transactionService) *TransactionHandler {
	return &TransactionHandler{infoService}
}

func (handler *TransactionHandler) Routes(c *gin.RouterGroup) {
	c.GET("/info", handler.GetInfo)
	c.POST("/sendCoin", handler.SendCoin)
	c.GET("/buy/:item", handler.BuyItem)
}

func (h TransactionHandler) GetInfo(c *gin.Context) {
	claims, _ := middleware.GetUser(c)
	response, err := h.transactionService.GetInfo(claims.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h TransactionHandler) SendCoin(c *gin.Context) {
	var request model.SendCoinRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{"Required fields are empty or not valid"})
		return
	}

	claims, _ := middleware.GetUser(c)
	err := h.transactionService.SendCoin(claims.UserId, request.ToUser, request.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h TransactionHandler) BuyItem(c *gin.Context) {
	item := c.Param("item")
	if item == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{"item name is not provided"})
		return
	}
	claims, _ := middleware.GetUser(c)
	err := h.transactionService.BuyItem(claims.UserId, item)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
