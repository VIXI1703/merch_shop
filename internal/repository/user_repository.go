package repository

import (
	"errors"
	"gorm.io/gorm"
	"merch_shop/internal/entity"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{
		db: db,
	}
}

func (repo *GormUserRepository) FindUserById(userId uint) (*entity.User, error) {
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

func (repo *GormUserRepository) FindUserByName(name string) (*entity.User, error) {
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

func (repo *GormUserRepository) UpdateUser(user *entity.User) error {
	return repo.db.Save(user).Error
}

func (repo *GormUserRepository) CreateUser(user *entity.User) error {
	return repo.db.Create(user).Error
}
