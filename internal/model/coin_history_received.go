package model

type CoinHistoryReceived struct {
	FromUser string `json:"fromUser"`
	Amount   uint   `json:"amount"`
}
