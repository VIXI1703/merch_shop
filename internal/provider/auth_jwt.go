package provider

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UserClaims struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTAuth struct {
	signingKey []byte
	expiration time.Duration
}

func NewJWTAuth(signingKey []byte, expiration time.Duration) *JWTAuth {
	return &JWTAuth{signingKey: signingKey, expiration: expiration}
}

func (auth JWTAuth) VerifyToken(tokenString string) (UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return auth.signingKey, nil
	})
	if err != nil {
		return UserClaims{}, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return *claims, nil
	}

	return UserClaims{}, errors.New("invalid token")
}

func (auth JWTAuth) GenerateToken(userId uint) (string, error) {
	claims := UserClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(auth.expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(auth.signingKey)
}
