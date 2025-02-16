package repository

import (
	"errors"
	"gorm.io/gorm"
	"merch_shop/internal/entity"
)

type UserRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) FindUserById(userId uint) (*entity.User, error) {
	user := new(entity.User)
	err := repo.db.First(user, userId).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) FindUserByName(name string) (*entity.User, error) {
	user := new(entity.User)
	err := repo.db.Where("name = ?", name).First(user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) CreateUser(user *entity.User) error {
	return repo.db.Create(user).Error
}
