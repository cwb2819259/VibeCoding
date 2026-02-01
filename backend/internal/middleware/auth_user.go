package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const UserIDKey = "user_id"

type UserClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthUser(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization"})
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization"})
			return
		}
		token, err := jwt.ParseWithClaims(parts[1], &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}
		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

// OptionalAuthUser 可选鉴权：有有效 token 时设置 user_id，无 token 或无效时设 0 并继续
func OptionalAuthUser(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.Set(UserIDKey, uint64(0))
			c.Next()
			return
		}
		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Set(UserIDKey, uint64(0))
			c.Next()
			return
		}
		token, err := jwt.ParseWithClaims(parts[1], &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.Set(UserIDKey, uint64(0))
			c.Next()
			return
		}
		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			c.Set(UserIDKey, uint64(0))
			c.Next()
			return
		}
		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint64 {
	v, _ := c.Get(UserIDKey)
	if v == nil {
		return 0
	}
	id, _ := v.(uint64)
	return id
}
