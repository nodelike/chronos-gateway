package handlers

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nodelike/chronos-gateway/internal/services"
)

func MediaUploadHandler(minio *services.MinIOClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user identification from headers or query
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = c.Query("user_id")
			if userID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user ID required"})
				return
			}
		}

		// Get device ID
		deviceID := c.GetHeader("X-Device-ID")
		if deviceID == "" {
			deviceID = c.Query("device_id")
			if deviceID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "device ID required"})
				return
			}
		}

		// Get file from form
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
			return
		}

		// Open the file
		fileContent, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
			return
		}
		defer fileContent.Close()

		// Read file content
		buffer := make([]byte, file.Size)
		_, err = fileContent.Read(buffer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
			return
		}

		// Generate object name
		contentType := file.Header.Get("Content-Type")
		extension := filepath.Ext(file.Filename)
		if extension == "" {
			extension = ".bin"
		}

		objectName := userID + "/" + deviceID + "/" + time.Now().Format("20060102-150405") + extension

		// Upload to MinIO
		err = minio.UploadFile(objectName, buffer, contentType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "uploaded",
			"file":   file.Filename,
			"path":   objectName,
		})
	}
}
