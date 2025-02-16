package model

type InfoResponse struct {
	Coins uint `json:"coins"`

	Inventory []Inventory `json:"inventory"`

	CoinHistory CoinHistory `json:"coinHistory"`
}
