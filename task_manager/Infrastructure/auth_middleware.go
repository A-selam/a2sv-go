package infrastructure

import (
	"strings"
	domain "task_manager/Domain"

	"github.com/gin-gonic/gin"
)

func NewAuthMiddleware(jwtService domain.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		username, role, err := jwtService.ParseToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("username", username)
		c.Set("role", role)
		c.Next()
	}
}
