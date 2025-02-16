package db

import (
	"gorm.io/gorm"
	"merch_shop/internal/entity"
)

func InitDB(gormDB *gorm.DB) *gorm.DB {
	err := gormDB.AutoMigrate(entity.InventoryItem{}, entity.Item{}, entity.User{}, entity.Transaction{})
	if err != nil {
		return nil
	}

	gormDB = SeedData(gormDB)

	return gormDB
}

func SeedData(db *gorm.DB) *gorm.DB {
	var count int64
	db.Model(&entity.Item{}).Count(&count)
	if count == 0 {
		items := []entity.Item{
			{
				Name:  "t-shirt",
				Price: 80,
			},
			{
				Name:  "cup",
				Price: 20,
			},
			{
				Name:  "book",
				Price: 50,
			},
			{
				Name:  "pen",
				Price: 10,
			},
			{
				Name:  "powerbank",
				Price: 200,
			},
			{
				Name:  "hoody",
				Price: 300,
			},
			{
				Name:  "umbrella",
				Price: 200,
			},
			{
				Name:  "socks",
				Price: 10,
			},
			{
				Name:  "wallet",
				Price: 50,
			},
			{
				Name:  "pink-hoody",
				Price: 500,
			},
		}

		if db.Create(&items).Error != nil {
			return nil
		}
	}

	return db
}
