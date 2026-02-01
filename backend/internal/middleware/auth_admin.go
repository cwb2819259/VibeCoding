package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const AdminIDKey = "admin_id"

type AdminClaims struct {
	AdminID uint64 `json:"admin_id"`
	jwt.RegisteredClaims
}

func AuthAdmin(secret string) gin.HandlerFunc {
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
		token, err := jwt.ParseWithClaims(parts[1], &AdminClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, ok := token.Claims.(*AdminClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}
		c.Set(AdminIDKey, claims.AdminID)
		c.Next()
	}
}

func GetAdminID(c *gin.Context) uint64 {
	v, _ := c.Get(AdminIDKey)
	if v == nil {
		return 0
	}
	id, _ := v.(uint64)
	return id
}
