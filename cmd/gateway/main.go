package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/nodelike/chronos-gateway/internal/config"
	"github.com/nodelike/chronos-gateway/internal/handlers"
	"github.com/nodelike/chronos-gateway/internal/middleware"
	"github.com/nodelike/chronos-gateway/internal/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	log.Println("Starting Chronos Gateway...")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize services
	log.Println("Initializing services...")
	kafkaProducer := services.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.DevelopmentMode)
	minioClient := services.NewMinIOClient(services.MinIOConfig{
		Endpoint:         cfg.MinIO.Endpoint,
		AccessKey:        cfg.MinIO.AccessKey,
		SecretKey:        cfg.MinIO.SecretKey,
		Bucket:           cfg.MinIO.Bucket,
		DevelopmentMode:  cfg.MinIO.DevelopmentMode,
		LocalStoragePath: cfg.MinIO.LocalStoragePath,
	})
	metrics := services.NewMetricsCollector()

	// Create Gin engine
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Public endpoints that don't require authentication
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Metrics endpoint (consider adding authentication for production)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes with authentication
	api := router.Group("")
	{
		// Apply authentication and metrics middleware to API routes
		api.Use(middleware.Metrics(metrics))
		api.Use(middleware.Authentication(cfg.APIKeys, cfg.DisableAuth))

		// API v1 routes
		v1 := api.Group("/v1")
		{
			// HTTP endpoints
			v1.POST("/android", handlers.HandleAndroidEvent(kafkaProducer))
			v1.POST("/macos", handlers.HandleMacOSEvent(kafkaProducer))
			v1.POST("/browser", handlers.HandleBrowserEvent(kafkaProducer, minioClient))

			// Location endpoints
			v1.POST("/location", handlers.HandleLocationEvent(kafkaProducer))
			v1.POST("/locations/batch", handlers.HandleBatchLocationEvents(kafkaProducer))

			// WebSocket endpoint
			v1.GET("/ws", handlers.WebSocketHandler(kafkaProducer))

			// Media upload endpoint
			v1.POST("/upload", handlers.MediaUploadHandler(minioClient))
		}
	}

	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start gRPC server in goroutine
	log.Println("Starting gRPC server on", cfg.GRPC.Port)
	go handlers.StartGRPCServer(cfg.GRPC.Port, kafkaProducer)

	// Start HTTP server in a goroutine
	log.Println("Starting HTTP server on", cfg.HTTP.Port)
	go func() {
		if err := router.Run(cfg.HTTP.Port); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	log.Println("Chronos Gateway is ready to receive events")
	log.Println("API Endpoints:")
	log.Printf("  - Health Check: http://localhost%s/health", cfg.HTTP.Port)
	log.Printf("  - Metrics: http://localhost%s/metrics", cfg.HTTP.Port)
	log.Printf("  - Location Events: http://localhost%s/v1/location", cfg.HTTP.Port)

	// Block until we receive a signal
	<-quit
	log.Println("Shutting down servers...")
}
