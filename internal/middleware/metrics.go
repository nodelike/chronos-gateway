package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nodelike/chronos-gateway/internal/services"
)

// Key for storing metrics in context
const MetricsKey = "metrics_collector"

func Metrics(metrics *services.MetricsCollector) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Store metrics collector in context for handlers to use
		c.Set(MetricsKey, metrics)

		// Process request
		c.Next()

		// Record metrics after processing
		status := strconv.Itoa(c.Writer.Status())
		latency := time.Since(start).Seconds()

		// Increment request counter with labels
		metrics.RequestCounter.WithLabelValues(method, path, status).Inc()

		// Record request duration
		metrics.RequestDuration.WithLabelValues(method, path).Observe(latency)
	}
}

// GetMetricsFromContext retrieves the metrics collector from the gin context
func GetMetricsFromContext(c *gin.Context) *services.MetricsCollector {
	if mc, exists := c.Get(MetricsKey); exists {
		if metrics, ok := mc.(*services.MetricsCollector); ok {
			return metrics
		}
	}
	return nil
}
