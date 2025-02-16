package service

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"merch_shop/internal/entity"
	"merch_shop/internal/model"
	"merch_shop/internal/repository"
	"testing"
)

// Mock repositories and unit of work

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByName(name string) (*entity.User, error) {
	args := m.Called(name)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindUserById(userId uint) (*entity.User, error) {
	args := m.Called(userId)
	return args.Get(0).(*entity.User), args.Error(1)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) AddItem(userId uint, itemId uint) error {
	args := m.Called(userId, itemId)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetItemByName(name string) (*entity.Item, error) {
	args := m.Called(name)
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockTransactionRepository) CreateTransaction(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetOutcomeTransactions(userId uint) ([]entity.Transaction, error) {
	args := m.Called(userId)
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetIncomeTransactions(userId uint) ([]entity.Transaction, error) {
	args := m.Called(userId)
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetUserInventory(userId uint) ([]entity.InventoryItem, error) {
	args := m.Called(userId)
	return args.Get(0).([]entity.InventoryItem), args.Error(1)
}

type MockTransactionUnitOfWork struct {
	UserRepo        *MockUserRepository
	TransactionRepo *MockTransactionRepository
	commitCalled    bool
	rollbackCalled  bool
}

func (m *MockTransactionUnitOfWork) BeginTransaction(opts ...*sql.TxOptions) (repository.TransactionUnitOfWork, error) {
	return m, nil
}

func (m *MockTransactionUnitOfWork) Commit() error {
	m.commitCalled = true
	return nil
}

func (m *MockTransactionUnitOfWork) Rollback() error {
	m.rollbackCalled = true
	return nil
}

func (m *MockTransactionUnitOfWork) UserRepository() repository.UserRepository {
	return m.UserRepo
}

func (m *MockTransactionUnitOfWork) TransactionRepository() repository.TransactionRepository {
	return m.TransactionRepo
}

type MockUnitOfWork struct {
	transactionUnitOfWork *MockTransactionUnitOfWork
}

func (m *MockUnitOfWork) BeginTransaction(opts ...*sql.TxOptions) (repository.TransactionUnitOfWork, error) {
	return m.transactionUnitOfWork, nil
}

func (m *MockUnitOfWork) UserRepository() repository.UserRepository {
	return m.transactionUnitOfWork.UserRepo
}

func (m *MockUnitOfWork) TransactionRepository() repository.TransactionRepository {
	return m.transactionUnitOfWork.TransactionRepo
}

// Tests

func TestTransactionService_GetInfo(t *testing.T) {
	t.Run("UserNotFound", func(t *testing.T) {
		userRepo := &MockUserRepository{}
		userRepo.On("FindUserById", uint(1)).Return((*entity.User)(nil), nil)

		tuow := &MockTransactionUnitOfWork{
			UserRepo:        userRepo,
			TransactionRepo: &MockTransactionRepository{},
		}
		uow := &MockUnitOfWork{transactionUnitOfWork: tuow}
		service := NewTransactionService(uow)

		_, err := service.GetInfo(1)
		assert.EqualError(t, err, "user not found")
	})

	t.Run("Success", func(t *testing.T) {
		user := &entity.User{Model: gorm.Model{ID: 1}, Balance: 1000}
		outcome := []entity.Transaction{
			{ToUser: entity.User{Name: "user2"}, Amount: 100},
		}
		income := []entity.Transaction{
			{FromUser: entity.User{Name: "user3"}, Amount: 200},
		}
		inventory := []entity.InventoryItem{
			{Item: entity.Item{Name: "item1"}, Quantity: 2},
		}

		userRepo := &MockUserRepository{}
		userRepo.On("FindUserById", uint(1)).Return(user, nil)

		transactionRepo := &MockTransactionRepository{}
		transactionRepo.On("GetOutcomeTransactions", uint(1)).Return(outcome, nil)
		transactionRepo.On("GetIncomeTransactions", uint(1)).Return(income, nil)
		transactionRepo.On("GetUserInventory", uint(1)).Return(inventory, nil)

		tuow := &MockTransactionUnitOfWork{
			UserRepo:        userRepo,
			TransactionRepo: transactionRepo,
		}
		uow := &MockUnitOfWork{transactionUnitOfWork: tuow}
		service := NewTransactionService(uow)

		res, err := service.GetInfo(1)
		assert.NoError(t, err)
		assert.Equal(t, model.InfoResponse{
			Coins: 1000,
			Inventory: []model.Inventory{
				{Name: "item1", Quantity: 2},
			},
			CoinHistory: model.CoinHistory{
				Sent: []model.CoinHistorySent{
					{ToUser: "user2", Amount: 100},
				},
				Received: []model.CoinHistoryReceived{
					{FromUser: "user3", Amount: 200},
				},
			},
		}, res)
	})
}

func TestTransactionService_SendCoin(t *testing.T) {
	t.Run("InsufficientBalance", func(t *testing.T) {
		fromUser := &entity.User{Model: gorm.Model{ID: 1}, Balance: 50}
		toUser := &entity.User{Model: gorm.Model{ID: 2}, Name: "user2"}

		userRepo := &MockUserRepository{}
		userRepo.On("FindUserById", uint(1)).Return(fromUser, nil)
		userRepo.On("FindUserByName", "user2").Return(toUser, nil)

		tuow := &MockTransactionUnitOfWork{
			UserRepo:        userRepo,
			TransactionRepo: &MockTransactionRepository{},
		}
		uow := &MockUnitOfWork{transactionUnitOfWork: tuow}
		service := NewTransactionService(uow)

		err := service.SendCoin(1, "user2", 100)
		assert.EqualError(t, err, "insufficient balance")
		assert.True(t, tuow.rollbackCalled)
	})

	t.Run("Success", func(t *testing.T) {
		fromUser := &entity.User{Model: gorm.Model{ID: 1}, Balance: 200}
		toUser := &entity.User{Model: gorm.Model{ID: 2}, Name: "user2", Balance: 0}

		userRepo := &MockUserRepository{}
		userRepo.On("FindUserById", uint(1)).Return(fromUser, nil)
		userRepo.On("FindUserByName", "user2").Return(toUser, nil)
		userRepo.On("UpdateUser", mock.Anything).Return(nil)

		transactionRepo := &MockTransactionRepository{}
		transactionRepo.On("CreateTransaction", mock.Anything).Return(nil)

		tuow := &MockTransactionUnitOfWork{
			UserRepo:        userRepo,
			TransactionRepo: transactionRepo,
		}
		uow := &MockUnitOfWork{transactionUnitOfWork: tuow}
		service := NewTransactionService(uow)

		err := service.SendCoin(1, "user2", 100)
		assert.NoError(t, err)
		assert.Equal(t, uint(100), fromUser.Balance)
		assert.Equal(t, uint(100), toUser.Balance)
		assert.True(t, tuow.commitCalled)
	})
}

func TestTransactionService_BuyItem(t *testing.T) {
	t.Run("ItemNotFound", func(t *testing.T) {
		user := &entity.User{Model: gorm.Model{ID: 1}, Balance: 1000}
		item := (*entity.Item)(nil)

		userRepo := &MockUserRepository{}
		userRepo.On("FindUserById", uint(1)).Return(user, nil)

		transactionRepo := &MockTransactionRepository{}
		transactionRepo.On("GetItemByName", "item1").Return(item, nil)

		tuow := &MockTransactionUnitOfWork{
			UserRepo:        userRepo,
			TransactionRepo: transactionRepo,
		}
		uow := &MockUnitOfWork{transactionUnitOfWork: tuow}
		service := NewTransactionService(uow)

		err := service.BuyItem(1, "item1")
		assert.EqualError(t, err, "item not found")
		assert.True(t, tuow.rollbackCalled)
	})

	t.Run("Success", func(t *testing.T) {
		user := &entity.User{Model: gorm.Model{ID: 1}, Balance: 1000}
		item := &entity.Item{Model: gorm.Model{ID: 1}, Name: "item1", Price: 500}

		userRepo := &MockUserRepository{}
		userRepo.On("FindUserById", uint(1)).Return(user, nil)
		userRepo.On("UpdateUser", mock.Anything).Return(nil)

		transactionRepo := &MockTransactionRepository{}
		transactionRepo.On("GetItemByName", "item1").Return(item, nil)
		transactionRepo.On("AddItem", uint(1), uint(1)).Return(nil)

		tuow := &MockTransactionUnitOfWork{
			UserRepo:        userRepo,
			TransactionRepo: transactionRepo,
		}
		uow := &MockUnitOfWork{transactionUnitOfWork: tuow}
		service := NewTransactionService(uow)

		err := service.BuyItem(1, "item1")
		assert.NoError(t, err)
		assert.Equal(t, uint(500), user.Balance)
		assert.True(t, tuow.commitCalled)
	})
}
