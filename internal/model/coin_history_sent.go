package model

type CoinHistorySent struct {
	ToUser string `json:"toUser"`
	Amount uint   `json:"amount"`
}
