package model

type CoinHistory struct {
	Received []CoinHistoryReceived `json:"received"`
	Sent     []CoinHistorySent     `json:"sent"`
}
