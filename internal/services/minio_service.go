package services

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Type alias to match the config struct
type MinIOConfig struct {
	Endpoint         string
	AccessKey        string
	SecretKey        string
	Bucket           string
	DevelopmentMode  bool
	LocalStoragePath string
}

type MinIOClient struct {
	client           *minio.Client
	bucket           string
	developmentMode  bool
	localStoragePath string
}

func NewMinIOClient(config MinIOConfig) *MinIOClient {
	// If in development mode, we don't connect to MinIO
	if config.DevelopmentMode {
		log.Println("Starting MinIO client in development mode - files will be saved to disk at", config.LocalStoragePath)
		// Create local storage path if it doesn't exist
		if config.LocalStoragePath == "" {
			config.LocalStoragePath = "./storage"
		}

		if err := os.MkdirAll(config.LocalStoragePath, 0755); err != nil {
			log.Printf("Error creating local storage path: %v", err)
		}

		return &MinIOClient{
			client:           nil,
			bucket:           config.Bucket,
			developmentMode:  true,
			localStoragePath: config.LocalStoragePath,
		}
	}

	// Connect to MinIO in production mode
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: true,
	})

	if err != nil {
		log.Printf("Error connecting to MinIO: %v", err)
		log.Println("Falling back to development mode - files will be saved to disk")

		// Create local storage path if it doesn't exist
		if config.LocalStoragePath == "" {
			config.LocalStoragePath = "./storage"
		}

		if err := os.MkdirAll(config.LocalStoragePath, 0755); err != nil {
			log.Printf("Error creating local storage path: %v", err)
		}

		return &MinIOClient{
			client:           nil,
			bucket:           config.Bucket,
			developmentMode:  true,
			localStoragePath: config.LocalStoragePath,
		}
	}

	return &MinIOClient{
		client:           client,
		bucket:           config.Bucket,
		developmentMode:  false,
		localStoragePath: config.LocalStoragePath,
	}
}

func (m *MinIOClient) UploadFile(objectName string, data []byte, contentType string) error {
	// In development mode, save to local file
	if m.developmentMode {
		// Create path if it doesn't exist
		fullPath := filepath.Join(m.localStoragePath, filepath.Dir(objectName))
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("error creating directory: %w", err)
		}

		// Save file
		filePath := filepath.Join(m.localStoragePath, objectName)
		if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}

		log.Printf("[DEV MODE] Saved file to %s (%d bytes)", filePath, len(data))
		return nil
	}

	// Upload to MinIO in production mode
	ctx := context.Background()
	_, err := m.client.PutObject(ctx, m.bucket, objectName,
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType})
	return err
}
