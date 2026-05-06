package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DatabaseURL string

	AWSRegion       string
	AWSAccessKey    string
	AWSSecretKey    string
	AWSEndpointURL  string

	SNSProductTopicARN          string
	SQSProductCreatedQueueURL   string

	CatalogGatewayBaseURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:                   getEnv("APP_PORT", "8080"),
		DatabaseURL:               os.Getenv("DATABASE_URL"),
		AWSRegion:                 getEnv("AWS_REGION", "us-east-1"),
		AWSAccessKey:              os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretKey:              os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSEndpointURL:            os.Getenv("AWS_ENDPOINT_URL"),
		SNSProductTopicARN:        os.Getenv("SNS_PRODUCT_TOPIC_ARN"),
		SQSProductCreatedQueueURL: os.Getenv("SQS_PRODUCT_CREATED_QUEUE_URL"),
		CatalogGatewayBaseURL:     os.Getenv("CATALOG_GATEWAY_BASE_URL"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.SNSProductTopicARN == "" {
		return nil, fmt.Errorf("SNS_PRODUCT_TOPIC_ARN is required")
	}
	if cfg.SQSProductCreatedQueueURL == "" {
		return nil, fmt.Errorf("SQS_PRODUCT_CREATED_QUEUE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
