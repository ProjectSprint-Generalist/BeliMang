package config

import (
	"fmt"
	"os"
)

type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	JWTSecret   string
	MinIO       MinIOConfig
}

type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

func LoadConfig() *Config {

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	cfg := &Config{
		Environment: os.Getenv("ENVIRONMENT"),
		Port:        os.Getenv("PORT"),
		DatabaseURL: dbURL,
		JWTSecret:   os.Getenv("JWT_SECRET"),
		MinIO: MinIOConfig{
			Endpoint:        os.Getenv("MINIO_ENDPOINT"),
			AccessKeyID:     os.Getenv("MINIO_ACCESS_KEY"),
			SecretAccessKey: os.Getenv("MINIO_SECRET_KEY"),
			UseSSL:          os.Getenv("MINIO_USE_SSL") == "true",
			BucketName:      os.Getenv("MINIO_BUCKET_NAME"),
		},
	}

	return cfg
}
