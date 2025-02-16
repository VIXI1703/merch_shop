package integration_test

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"merch_shop/internal/db"
	"merch_shop/internal/entity"
	"merch_shop/internal/provider"
	"merch_shop/internal/repository"
	"merch_shop/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	startBalance = uint(1000)
	jwtSecret    = "test_secret"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dbConn, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{SkipDefaultTransaction: true})
	assert.NoError(t, err)

	db.InitDB(dbConn)
	return dbConn
}

func TestAuthServiceIntegration(t *testing.T) {
	db := setupTestDB(t)
	uow := repository.NewGormUnitOfWork(db)
	jwtAuth := provider.NewJWTAuth([]byte(jwtSecret), 24*time.Hour)
	authService := service.NewAuthService(jwtAuth, uow)

	t.Run("NewUserRegistration", func(t *testing.T) {
		token, err := authService.Authenticate("newuser", "password")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify user creation
		userRepo := uow.UserRepository()
		user, err := userRepo.FindUserByName("newuser")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, startBalance, user.Balance)
	})

	t.Run("ExistingUserLogin", func(t *testing.T) {
		// Create user first
		_, _ = authService.Authenticate("existinguser", "password")

		t.Run("ValidCredentials", func(t *testing.T) {
			token, err := authService.Authenticate("existinguser", "password")
			assert.NoError(t, err)
			assert.NotEmpty(t, token)
		})

		t.Run("InvalidCredentials", func(t *testing.T) {
			_, err := authService.Authenticate("existinguser", "wrongpass")
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "password is incorrect")
		})
	})
}

func TestTransactionServiceIntegration(t *testing.T) {
	db := setupTestDB(t)
	uow := repository.NewGormUnitOfWork(db)
	service := service.NewTransactionService(uow)

	// Create test users
	userRepo := uow.UserRepository()
	user1 := createTestUser(t, uow, "user1", startBalance)
	user2 := createTestUser(t, uow, "user2", startBalance)

	t.Run("GetUserInfo", func(t *testing.T) {
		info, err := service.GetInfo(user1.ID)
		assert.NoError(t, err)
		assert.Equal(t, startBalance, info.Coins)
		assert.Empty(t, info.Inventory)
	})

	t.Run("CoinTransfer", func(t *testing.T) {
		transferAmount := uint(200)

		// Initial balance verification
		user1Before, _ := userRepo.FindUserById(user1.ID)
		user2Before, _ := userRepo.FindUserById(user2.ID)

		// Perform transfer
		err := service.SendCoin(user1.ID, user2.Name, transferAmount)
		assert.NoError(t, err)

		// Verify balances
		user1After, _ := userRepo.FindUserById(user1.ID)
		user2After, _ := userRepo.FindUserById(user2.ID)

		assert.Equal(t, user1Before.Balance-transferAmount, user1After.Balance)
		assert.Equal(t, user2Before.Balance+transferAmount, user2After.Balance)

		// Verify transaction history
		txRepo := uow.TransactionRepository()
		transactions, _ := txRepo.GetOutcomeTransactions(user1.ID)
		assert.Len(t, transactions, 1)
		assert.Equal(t, user2.ID, transactions[0].ToId)
	})

	t.Run("ItemPurchase", func(t *testing.T) {
		itemName := "t-shirt"
		user, _ := userRepo.FindUserById(user1.ID)
		initialBalance := user.Balance

		// Get item price
		txRepo := uow.TransactionRepository()
		item, _ := txRepo.GetItemByName(itemName)

		err := service.BuyItem(user1.ID, itemName)
		assert.NoError(t, err)

		// Verify balance deduction
		userAfter, _ := userRepo.FindUserById(user1.ID)
		assert.Equal(t, initialBalance-item.Price, userAfter.Balance)

		// Verify inventory
		inventory, _ := txRepo.GetUserInventory(user1.ID)
		assert.Len(t, inventory, 1)
		assert.Equal(t, itemName, inventory[0].Item.Name)
		assert.Equal(t, uint(1), inventory[0].Quantity)
	})
}

func TestTransactionServiceEdgeCases(t *testing.T) {
	db := setupTestDB(t)
	uow := repository.NewGormUnitOfWork(db)
	service := service.NewTransactionService(uow)

	t.Run("SendToNonExistentUser", func(t *testing.T) {
		sender := createTestUser(t, uow, "sender1", 1000)
		err := service.SendCoin(sender.ID, "ghost_user", 100)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("OverdraftPrevention", func(t *testing.T) {
		sender := createTestUser(t, uow, "sender2", 100)
		receiver := createTestUser(t, uow, "receiver2", 0)

		err := service.SendCoin(sender.ID, receiver.Name, 200)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient balance")
	})

	t.Run("BuyNonExistentItem", func(t *testing.T) {
		user := createTestUser(t, uow, "buyer1", 1000)
		err := service.BuyItem(user.ID, "unicorn")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "item not found")
	})

	t.Run("SelfTransferPrevention", func(t *testing.T) {
		user := createTestUser(t, uow, "selfsender", 1000)
		err := service.SendCoin(user.ID, user.Name, 100)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot send coin to yourself")
	})
}

func createTestUser(t *testing.T, uow repository.UnitOfWork, name string, balance uint) *entity.User {
	user := &entity.User{
		Name:         name,
		PasswordHash: "hash",
		Balance:      balance,
	}
	err := uow.UserRepository().CreateUser(user)
	assert.NoError(t, err)
	return user
}
