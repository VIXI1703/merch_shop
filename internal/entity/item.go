package entity

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Name  string `gorm:"uniqueIndex:item_name"`
	Price uint
}
