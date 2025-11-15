package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const hardcodedAPIKey = "my-secret-key-123"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")

		if apiKey != hardcodedAPIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid API Key"})
			return
		}

		c.Next()
	}
}
