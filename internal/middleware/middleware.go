package middleware

import (
	"github.com/gin-gonic/gin"
	"merch_shop/internal/model"
	"merch_shop/internal/provider"
	"net/http"
	"strings"
)

func JWTAuthMiddleware(auth *provider.JWTAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenWithPrefix := c.GetHeader("Authorization")
		token, found := strings.CutPrefix(tokenWithPrefix, "Bearer ")
		if !found || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Errors: "Authorization header required"})
			return
		}
		claims, err := auth.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Errors: "Token invalid"})
			return
		}
		c.Set("user", claims)
		c.Next()
	}
}

func GetUser(c *gin.Context) (provider.UserClaims, bool) {
	userClaims, found := c.Get("user")
	if !found {
		return provider.UserClaims{}, false
	}
	return userClaims.(provider.UserClaims), found
}
