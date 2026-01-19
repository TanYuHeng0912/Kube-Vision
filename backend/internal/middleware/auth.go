package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates authentication tokens
func AuthMiddleware(authEnabled bool, authToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth if disabled
		if !authEnabled {
			c.Next()
			return
		}

		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token (format: "Bearer <token>" or just "<token>")
		token := strings.TrimSpace(authHeader)
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		// Validate token
		if token != authToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

