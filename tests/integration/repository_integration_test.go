package integration_test

import (
	"merch_shop/internal/db"
	"testing"

	"merch_shop/internal/entity"
	"merch_shop/internal/repository"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupRepoTestDB(t *testing.T) *gorm.DB {
	dbConn, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{SkipDefaultTransaction: true})
	assert.NoError(t, err)

	dbConn = db.InitDB(dbConn)
	return dbConn
}

func TestUserRepositoryIntegration(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := repository.NewGormUserRepository(db)

	t.Run("CreateDuplicateUser", func(t *testing.T) {
		user1 := &entity.User{Name: "alice", PasswordHash: "hash1"}
		assert.NoError(t, repo.CreateUser(user1))

		user2 := &entity.User{Name: "alice", PasswordHash: "hash2"}
		err := repo.CreateUser(user2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "UNIQUE constraint failed")
	})
}

func TestTransactionRepositoryIntegration(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := repository.NewGormTransactionRepository(db)

	t.Run("TransactionHistoryJoins", func(t *testing.T) {
		// Create users
		sender := &entity.User{Name: "sender", Balance: 1000}
		receiver := &entity.User{Name: "receiver", Balance: 1000}
		db.Create(sender)
		db.Create(receiver)

		// Create transaction
		tx := &entity.Transaction{
			FromId: sender.ID,
			ToId:   receiver.ID,
			Amount: 200,
		}
		db.Create(tx)

		t.Run("OutcomeTransactions", func(t *testing.T) {
			transactions, _ := repo.GetOutcomeTransactions(sender.ID)
			assert.Len(t, transactions, 1)
			assert.Equal(t, "receiver", transactions[0].ToUser.Name)
		})

		t.Run("IncomeTransactions", func(t *testing.T) {
			transactions, _ := repo.GetIncomeTransactions(receiver.ID)
			assert.Len(t, transactions, 1)
			assert.Equal(t, "sender", transactions[0].FromUser.Name)
		})
	})
}

func TestInventoryRepositoryIntegration(t *testing.T) {
	db := setupRepoTestDB(t)
	repo := repository.NewGormTransactionRepository(db)

	t.Run("UpsertInventoryItem", func(t *testing.T) {
		user := &entity.User{Name: "collector"}
		item := &entity.Item{Name: "sticker", Price: 5}
		db.Create(user)
		db.Create(item)

		// First addition
		assert.NoError(t, repo.AddItem(user.ID, item.ID))

		// Second addition
		assert.NoError(t, repo.AddItem(user.ID, item.ID))

		inventory, _ := repo.GetUserInventory(user.ID)
		assert.Len(t, inventory, 1)
		assert.Equal(t, uint(2), inventory[0].Quantity)
	})
}
