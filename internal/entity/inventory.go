package entity

import "gorm.io/gorm"

type InventoryItem struct {
	gorm.Model
	UserID   uint `gorm:"index:entry,unique"`
	User     User
	ItemID   uint `gorm:"index:entry,unique"`
	Item     Item
	Quantity uint
}
