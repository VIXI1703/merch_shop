package db

import (
	"fmt"
	_ "github.com/jackc/pgx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"merch_shop/internal/config"
)

func SetupDB(cfg *config.DB) *gorm.DB {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dataSourceName,
		PreferSimpleProtocol: true,
	}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		panic(err.Error())
	}

	sqlDB, _ := gormDB.DB()
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLife)

	return InitDB(gormDB)
}
