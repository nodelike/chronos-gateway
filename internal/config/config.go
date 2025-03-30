package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP struct {
		Port string `mapstructure:"Port"`
	} `mapstructure:"HTTP"`
	GRPC struct {
		Port string `mapstructure:"Port"`
	} `mapstructure:"GRPC"`
	Kafka struct {
		Brokers         []string `mapstructure:"Brokers"`
		DevelopmentMode bool     `mapstructure:"DevelopmentMode"`
	} `mapstructure:"Kafka"`
	MinIO struct {
		Endpoint         string `mapstructure:"Endpoint"`
		AccessKey        string `mapstructure:"AccessKey"`
		SecretKey        string `mapstructure:"SecretKey"`
		Bucket           string `mapstructure:"Bucket"`
		DevelopmentMode  bool   `mapstructure:"DevelopmentMode"`
		LocalStoragePath string `mapstructure:"LocalStoragePath"`
	} `mapstructure:"MinIO"`
	APIKeys     map[string]bool `mapstructure:"APIKeys"`
	DisableAuth bool            `mapstructure:"DisableAuth"`
	Metrics     struct {
		Enable   bool   `mapstructure:"Enable"`
		Endpoint string `mapstructure:"Endpoint"`
	} `mapstructure:"Metrics"`
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	// Log loaded configuration
	log.Println("Configuration loaded successfully")
	if cfg.Kafka.DevelopmentMode {
		log.Println("Kafka in development mode: messages will be logged")
	}
	if cfg.MinIO.DevelopmentMode {
		log.Println("MinIO in development mode: files will be saved to", cfg.MinIO.LocalStoragePath)
	}
	if cfg.DisableAuth {
		log.Println("API Authentication is disabled")
	}

	return &cfg
}
