package repository

import (
	"database/sql"
	"gorm.io/gorm"
	"merch_shop/internal/entity"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	FindUserByName(name string) (*entity.User, error)
	FindUserById(userId uint) (*entity.User, error)
}

type TransactionRepository interface {
	AddItem(userId uint, itemId uint) error
	GetItemByName(name string) (*entity.Item, error)
	CreateTransaction(transaction *entity.Transaction) error
	GetOutcomeTransactions(userId uint) ([]entity.Transaction, error)
	GetIncomeTransactions(userId uint) ([]entity.Transaction, error)
	GetUserInventory(userId uint) ([]entity.InventoryItem, error)
}

type UnitOfWork interface {
	BeginTransaction(opts ...*sql.TxOptions) (TransactionUnitOfWork, error)
	UserRepository() UserRepository
	TransactionRepository() TransactionRepository
}

type TransactionUnitOfWork interface {
	UnitOfWork
	Commit() error
	Rollback() error
}

type GormUnitOfWork struct {
	db *gorm.DB
}

func NewGormUnitOfWork(db *gorm.DB) *GormUnitOfWork {
	return &GormUnitOfWork{db: db}
}

func (u *GormUnitOfWork) BeginTransaction(opts ...*sql.TxOptions) (TransactionUnitOfWork, error) {
	db := u.db.Begin(opts...)
	if db.Error != nil {
		return nil, db.Error
	}

	return &GormUnitOfWork{
		db: db,
	}, nil
}

func (u *GormUnitOfWork) Commit() error {
	return u.db.Commit().Error
}

func (u *GormUnitOfWork) Rollback() error {
	return u.db.Rollback().Error
}

func (u *GormUnitOfWork) UserRepository() UserRepository {
	return NewGormUserRepository(u.db)
}

func (u *GormUnitOfWork) TransactionRepository() TransactionRepository {
	return NewGormTransactionRepository(u.db)
}
