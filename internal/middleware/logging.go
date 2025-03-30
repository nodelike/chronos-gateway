package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		// Log to console for now, could use a proper logger in production
		log := fmt.Sprintf("[GIN] %s | %d | %s | %s | %s | %s\n",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency.String(),
			clientIP,
			method,
			path)

		gin.DefaultWriter.Write([]byte(log))
	}
}
