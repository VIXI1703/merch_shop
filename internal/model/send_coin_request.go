package model

type SendCoinRequest struct {
	ToUser string `json:"toUser" binding:"required"`
	Amount uint   `json:"amount" binding:"required"`
}
