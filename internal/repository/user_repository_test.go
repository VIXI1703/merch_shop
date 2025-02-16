package repository

import (
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"merch_shop/internal/entity"
	"testing"
)

func setupUserDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&entity.User{})
	return db
}

func TestGormUserRepository_FindUserById(t *testing.T) {
	db := setupUserDB()
	repo := NewGormUserRepository(db)

	t.Run("UserExists", func(t *testing.T) {
		user := &entity.User{Name: "test", PasswordHash: "hash", Balance: 1000}
		db.Create(user)

		found, err := repo.FindUserById(user.ID)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		found, err := repo.FindUserById(999)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestGormUserRepository_FindUserByName(t *testing.T) {
	db := setupUserDB()
	repo := NewGormUserRepository(db)

	t.Run("UserExists", func(t *testing.T) {
		user := &entity.User{Name: "alice", PasswordHash: "hash", Balance: 1000}
		db.Create(user)

		found, err := repo.FindUserByName("alice")
		assert.NoError(t, err)
		assert.Equal(t, user.ID, found.ID)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		found, err := repo.FindUserByName("bob")
		assert.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestGormUserRepository_CreateUser(t *testing.T) {
	db := setupUserDB()
	repo := NewGormUserRepository(db)

	user := &entity.User{Name: "test", PasswordHash: "hash"}
	err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)

	var count int64
	db.Model(&entity.User{}).Where("name = ?", "test").Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestGormUserRepository_UpdateUser(t *testing.T) {
	db := setupUserDB()
	repo := NewGormUserRepository(db)

	user := &entity.User{Name: "test", Balance: 500}
	db.Create(user)

	user.Balance = 1000
	err := repo.UpdateUser(user)
	assert.NoError(t, err)

	var updatedUser entity.User
	db.First(&updatedUser, user.ID)
	assert.Equal(t, uint(1000), updatedUser.Balance)
}
