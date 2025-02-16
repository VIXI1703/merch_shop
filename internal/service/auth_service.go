package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"merch_shop/internal/entity"
	"merch_shop/internal/provider"
	"merch_shop/internal/repository"
)

const START_BALANCE = 1000

type AuthService struct {
	jwtAuth *provider.JWTAuth
	uow     repository.UnitOfWork
}

func NewAuthService(jwtAuth *provider.JWTAuth, uow repository.UnitOfWork) AuthService {
	return AuthService{jwtAuth: jwtAuth, uow: uow}
}

func (auth AuthService) Authenticate(username, password string) (string, error) {
	userRepository := auth.uow.UserRepository()

	user, err := userRepository.FindUserByName(username)
	if err != nil {
		return "", fmt.Errorf("failed to find user by username %s", username)
	}
	if user == nil {
		passwordHash, err := hashPassword(password)
		if err != nil {
			return "", fmt.Errorf("failed to hash password")
		}
		user = &entity.User{
			Name:         username,
			PasswordHash: passwordHash,
			Balance:      START_BALANCE,
		}

		if err := userRepository.CreateUser(user); err != nil {
			return "", fmt.Errorf("failed to create user")
		}
	} else if !verifyPassword(password, user.PasswordHash) {
		return "", fmt.Errorf("password is incorrect")
	}

	token, err := auth.jwtAuth.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token")
	}
	return token, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
