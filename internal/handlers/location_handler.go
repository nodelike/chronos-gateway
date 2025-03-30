package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nodelike/chronos-gateway/internal/middleware"
	"github.com/nodelike/chronos-gateway/internal/models"
	"github.com/nodelike/chronos-gateway/internal/services"
	"github.com/nodelike/chronos-gateway/internal/utils"
)

// HandleLocationEvent processes location data from React Native Background Geolocation
func HandleLocationEvent(producer *services.KafkaProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricsCollector := middleware.GetMetricsFromContext(c)

		var event models.LocationEvent
		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		receivedTime := time.Now()

		if event.Timestamp.IsZero() {
			event.Timestamp = receivedTime
		}

		latency := receivedTime.Sub(event.Timestamp).Seconds()

		if event.EventType == "" {
			event.EventType = "location" // Default event type
		}
		event.DeviceID = utils.SanitizeString(event.DeviceID)
		event.UserID = utils.SanitizeString(event.UserID)

		if metricsCollector != nil {
			metricsCollector.RecordLocationEvent(event.EventType, "android", latency)
		}
		producer.SendEvent("location", event.ToJSON())

		c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
	}
}

// HandleBatchLocationEvents processes batched location data
func HandleBatchLocationEvents(producer *services.KafkaProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get metrics collector from the context
		metricsCollector := middleware.GetMetricsFromContext(c)

		var events []models.LocationEvent
		if err := c.ShouldBindJSON(&events); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(events) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "empty batch"})
			return
		}

		// Record batch size
		if metricsCollector != nil {
			metricsCollector.RecordBatchSize("android", len(events))
		}

		// Record event receipt time for latency calculation
		receivedTime := time.Now()

		// Process each event in the batch
		for i := range events {
			// Set timestamp if not provided
			if events[i].Timestamp.IsZero() {
				events[i].Timestamp = receivedTime
			}

			// Calculate latency
			latency := receivedTime.Sub(events[i].Timestamp).Seconds()

			// Default event type if not provided
			if events[i].EventType == "" {
				events[i].EventType = "location"
			}

			// Sanitize fields
			events[i].DeviceID = utils.SanitizeString(events[i].DeviceID)
			events[i].UserID = utils.SanitizeString(events[i].UserID)

			// Record metrics for each event
			if metricsCollector != nil {
				metricsCollector.RecordLocationEvent(events[i].EventType, "android", latency)
			}

			// Send to Kafka
			producer.SendEvent("location", events[i].ToJSON())
		}

		c.JSON(http.StatusAccepted, gin.H{
			"status": "accepted",
			"count":  len(events),
		})
	}
}
