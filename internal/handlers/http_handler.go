package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nodelike/chronos-gateway/internal/models"
	"github.com/nodelike/chronos-gateway/internal/services"
)

func HandleAndroidEvent(producer *services.KafkaProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.AndroidEvent
		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set timestamp if not provided
		if event.Timestamp.IsZero() {
			event.Timestamp = time.Now()
		}

		// Send to Kafka
		producer.SendEvent("android", event.ToJSON())

		c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
	}
}

func HandleMacOSEvent(producer *services.KafkaProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.MacOSEvent
		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set timestamp if not provided
		if event.Timestamp.IsZero() {
			event.Timestamp = time.Now()
		}

		// Send to Kafka
		producer.SendEvent("macos", event.ToJSON())

		c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
	}
}

func HandleBrowserEvent(producer *services.KafkaProducer, minio *services.MinIOClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.BrowserEvent
		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set timestamp if not provided
		if event.Timestamp.IsZero() {
			event.Timestamp = time.Now()
		}

		// Process media files if present
		if event.HasMedia && len(event.Media) > 0 {
			objectName := event.UserID + "/" + event.DeviceID + "/" + time.Now().Format("20060102-150405") + ".bin"
			err := minio.UploadFile(objectName, event.Media, event.MediaType)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload media"})
				return
			}
		}

		// Send to Kafka
		producer.SendEvent("browser", event.ToJSON())

		c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
	}
}
