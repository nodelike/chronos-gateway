package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authentication(validKeys map[string]bool, disableAuth bool) gin.HandlerFunc {
	if disableAuth {
		log.Println("API Authentication disabled - all requests will be allowed")
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if !validKeys[apiKey] {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}
		c.Next()
	}
}
