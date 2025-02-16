package repository

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"merch_shop/internal/entity"
)

type GormTransactionRepository struct {
	db *gorm.DB
}

func NewGormTransactionRepository(db *gorm.DB) *GormTransactionRepository {
	return &GormTransactionRepository{
		db: db,
	}
}

func (repo *GormTransactionRepository) GetUserInventory(userId uint) ([]entity.InventoryItem, error) {
	var inventoryItems []entity.InventoryItem
	err := repo.db.Joins("Item").Where("user_id = ?", userId).Find(&inventoryItems).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.InventoryItem{}, nil
		}
		return nil, err
	}
	return inventoryItems, nil
}
func (repo *GormTransactionRepository) GetIncomeTransactions(userId uint) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := repo.db.Joins("FromUser").Where("to_id = ?", userId).Find(&transactions).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.Transaction{}, nil
		}
		return nil, err
	}
	return transactions, nil

}

func (repo *GormTransactionRepository) GetOutcomeTransactions(userId uint) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := repo.db.Joins("ToUser").Where("from_id = ?", userId).Find(&transactions).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.Transaction{}, nil
		}
		return nil, err
	}
	return transactions, nil
}

func (repo *GormTransactionRepository) CreateTransaction(transaction *entity.Transaction) error {
	return repo.db.Create(transaction).Error
}

func (repo *GormTransactionRepository) GetItemByName(name string) (*entity.Item, error) {
	item := new(entity.Item)
	err := repo.db.Where("name = ?", name).First(item).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (repo *GormTransactionRepository) AddItem(userId uint, itemId uint) error {
	return repo.db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "item_id"}}, // Составной ключ для поиска дубликатов
			DoUpdates: clause.Assignments(map[string]interface{}{
				"quantity": gorm.Expr("inventory_items.quantity + 1"), // Увеличиваем quantity на 1
			}),
		},
	).Create(&entity.InventoryItem{
		UserID:   userId,
		ItemID:   itemId,
		Quantity: 1,
	}).Error
}
