package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `gorm:"uniqueIndex:user_name"`
	PasswordHash string
	Balance      uint `gorm:"default:1000"`
}
