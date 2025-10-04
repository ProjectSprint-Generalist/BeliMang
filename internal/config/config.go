package config

import (
	"os"
)

type Config struct {
	Environment string
	Port        string
	JWTSecret   string
	DB          DBConfig
	MinIO       MinIOConfig
}

type MinIOConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
}

type DBConfig struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
	SSLMode  string
	URL      string
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadConfig() *Config {
	cfg := &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		DB:          *LoadDBConfig(),
		MinIO:       *LoadMinIOConfig(),
	}
	return cfg
}

func LoadDBConfig() *DBConfig {
	return &DBConfig{
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Name:     getEnv("DB_NAME", "belimang"),
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
		URL:      getEnv("DATABASE_URL", ""),
	}
}

func LoadMinIOConfig() *MinIOConfig {
	return &MinIOConfig{
		Endpoint:        getEnv("MINIO_ENDPOINT", "localhost:9000"),
		AccessKeyID:     getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		SecretAccessKey: getEnv("MINIO_SECRET_KEY", "minioadmin123"),
		UseSSL:          getEnv("MINIO_USE_SSL", "false") == "true",
		BucketName:      getEnv("MINIO_BUCKET_NAME", "belimang-files"),
	}
}
