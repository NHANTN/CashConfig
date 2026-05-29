package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/cashier-config/server/internal/service"
)

func JWTAuth(authSvc *service.AuthService, apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try API key first
		reqKey := c.GetHeader("X-API-Key")
		if reqKey != "" && apiKey != "" && reqKey == apiKey {
			c.Set("user_id", int64(0))
			c.Set("username", "api-key")
			c.Set("role_code", "api")
			c.Next()
			return
		}

		// Fall back to JWT
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "missing authorization header or X-API-Key"})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid authorization format"})
			return
		}
		claims, err := authSvc.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid or expired token"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role_code", claims.RoleCode)
		c.Next()
	}
}
