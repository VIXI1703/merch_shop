package repository

import (
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"merch_shop/internal/entity"
	"testing"
)

func setupTransactionDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&entity.User{}, &entity.Item{}, &entity.Transaction{}, &entity.InventoryItem{})
	return db
}

func TestGormTransactionRepository_CreateTransaction(t *testing.T) {
	db := setupTransactionDB()
	repo := NewGormTransactionRepository(db)

	fromUser := &entity.User{Name: "alice", Balance: 1000}
	toUser := &entity.User{Name: "bob", Balance: 500}
	db.Create(fromUser)
	db.Create(toUser)

	tx := &entity.Transaction{FromId: fromUser.ID, ToId: toUser.ID, Amount: 200}
	err := repo.CreateTransaction(tx)
	assert.NoError(t, err)
	assert.NotZero(t, tx.ID)
}

func TestGormTransactionRepository_GetItemByName(t *testing.T) {
	db := setupTransactionDB()
	repo := NewGormTransactionRepository(db)

	t.Run("ItemExists", func(t *testing.T) {
		item := &entity.Item{Name: "sword", Price: 100}
		db.Create(item)

		found, err := repo.GetItemByName("sword")
		assert.NoError(t, err)
		assert.Equal(t, item.ID, found.ID)
	})

	t.Run("ItemNotFound", func(t *testing.T) {
		found, err := repo.GetItemByName("shield")
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestGormTransactionRepository_AddItem(t *testing.T) {
	db := setupTransactionDB()
	repo := NewGormTransactionRepository(db)

	user := &entity.User{Name: "alice"}
	item := &entity.Item{Name: "potion", Price: 50}
	db.Create(user)
	db.Create(item)

	// First addition
	err := repo.AddItem(user.ID, item.ID)
	assert.NoError(t, err)

	var inventoryItem entity.InventoryItem
	db.Where("user_id = ? AND item_id = ?", user.ID, item.ID).First(&inventoryItem)
	assert.Equal(t, uint(1), inventoryItem.Quantity)

	// Second addition (upsert)
	err = repo.AddItem(user.ID, item.ID)
	assert.NoError(t, err)
	db.Where("user_id = ? AND item_id = ?", user.ID, item.ID).First(&inventoryItem)
	assert.Equal(t, uint(2), inventoryItem.Quantity)
}

func TestGormTransactionRepository_GetUserInventory(t *testing.T) {
	db := setupTransactionDB()
	repo := NewGormTransactionRepository(db)

	user := &entity.User{Name: "alice"}
	item := &entity.Item{Name: "armor", Price: 200}
	db.Create(user)
	db.Create(item)
	db.Create(&entity.InventoryItem{UserID: user.ID, ItemID: item.ID, Quantity: 3})

	inventory, err := repo.GetUserInventory(user.ID)
	assert.NoError(t, err)
	assert.Len(t, inventory, 1)
	assert.Equal(t, "armor", inventory[0].Item.Name)
	assert.Equal(t, uint(3), inventory[0].Quantity)
}

func TestGormTransactionRepository_GetIncomeTransactions(t *testing.T) {
	db := setupTransactionDB()
	repo := NewGormTransactionRepository(db)

	fromUser := &entity.User{Name: "alice"}
	toUser := &entity.User{Name: "bob"}
	db.Create(fromUser)
	db.Create(toUser)
	tx := &entity.Transaction{FromId: fromUser.ID, ToId: toUser.ID, Amount: 100}
	db.Create(tx)

	transactions, err := repo.GetIncomeTransactions(toUser.ID)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, "alice", transactions[0].FromUser.Name)
}

func TestGormTransactionRepository_GetOutcomeTransactions(t *testing.T) {
	db := setupTransactionDB()
	repo := NewGormTransactionRepository(db)

	fromUser := &entity.User{Name: "alice"}
	toUser := &entity.User{Name: "bob"}
	db.Create(fromUser)
	db.Create(toUser)
	tx := &entity.Transaction{FromId: fromUser.ID, ToId: toUser.ID, Amount: 100}
	db.Create(tx)

	transactions, err := repo.GetOutcomeTransactions(fromUser.ID)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, "bob", transactions[0].ToUser.Name)
}
