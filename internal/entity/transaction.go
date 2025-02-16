package entity

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	FromId   uint
	FromUser User `gorm:"foreignKey:FromId"`
	ToId     uint
	ToUser   User `gorm:"foreignKey:ToId"`
	Amount   uint
}
