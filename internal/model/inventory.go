package model

type Inventory struct {
	Name     string `json:"name,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}
