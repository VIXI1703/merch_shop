package service

import (
	"database/sql"
	"testing"
	"time"

	"merch_shop/internal/entity"
	"merch_shop/internal/provider"
	"merch_shop/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthUserRepository struct {
	mock.Mock
}

func (m *MockAuthUserRepository) CreateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthUserRepository) UpdateUser(user *entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockAuthUserRepository) FindUserByName(name string) (*entity.User, error) {
	args := m.Called(name)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthUserRepository) FindUserById(userId uint) (*entity.User, error) {
	args := m.Called(userId)
	return args.Get(0).(*entity.User), args.Error(1)
}

type MockAuthUnitOfWork struct {
	userRepo *MockAuthUserRepository
}

func (m *MockAuthUnitOfWork) BeginTransaction(opts ...*sql.TxOptions) (repository.TransactionUnitOfWork, error) {
	panic("not implemented")
}

func (m *MockAuthUnitOfWork) UserRepository() repository.UserRepository {
	return m.userRepo
}

func (m *MockAuthUnitOfWork) TransactionRepository() repository.TransactionRepository {
	panic("not implemented")
}

func TestAuthService_Authenticate(t *testing.T) {
	t.Run("NewUser", func(t *testing.T) {
		userRepo := &MockAuthUserRepository{}
		userRepo.On("FindUserByName", "newuser").Return((*entity.User)(nil), nil)
		userRepo.On("CreateUser", mock.AnythingOfType("*entity.User")).Return(nil)

		uow := &MockAuthUnitOfWork{userRepo: userRepo}
		jwtAuth := provider.NewJWTAuth([]byte("test_secret"), time.Hour)
		service := NewAuthService(jwtAuth, uow)

		token, err := service.Authenticate("newuser", "password")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		userRepo.AssertExpectations(t)
	})

	t.Run("ExistingUserCorrectPassword", func(t *testing.T) {
		hashedPassword, _ := hashPassword("password")
		existingUser := &entity.User{
			Name:         "existinguser",
			PasswordHash: hashedPassword,
		}

		userRepo := &MockAuthUserRepository{}
		userRepo.On("FindUserByName", "existinguser").Return(existingUser, nil)

		uow := &MockAuthUnitOfWork{userRepo: userRepo}
		jwtAuth := provider.NewJWTAuth([]byte("test_secret"), time.Hour)
		service := NewAuthService(jwtAuth, uow)

		token, err := service.Authenticate("existinguser", "password")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("ExistingUserIncorrectPassword", func(t *testing.T) {
		existingUser := &entity.User{
			Name:         "existinguser",
			PasswordHash: "wronghash",
		}

		userRepo := &MockAuthUserRepository{}
		userRepo.On("FindUserByName", "existinguser").Return(existingUser, nil)

		uow := &MockAuthUnitOfWork{userRepo: userRepo}
		jwtAuth := provider.NewJWTAuth([]byte("test_secret"), time.Hour)
		service := NewAuthService(jwtAuth, uow)

		_, err := service.Authenticate("existinguser", "password")
		assert.EqualError(t, err, "password is incorrect")
	})
}
