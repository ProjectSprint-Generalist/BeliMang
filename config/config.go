package config

import "github.com/ProjectSprint-Generalist/BeliMang/internal/shared"

type Config struct {
	Environment string
	Port        string
	JWT         JWTConfig
	DB          DBConfig
	MinIO       MinIOConfig
}

type JWTConfig struct {
	Secret string
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

func LoadConfig() *Config {
	cfg := &Config{
		Environment: shared.SafeGetEnv("ENVIRONMENT", "development"),
		Port:        shared.SafeGetEnv("PORT", "8080"),
		JWT:         *LoadJWTConfig(),
		DB:          *LoadDBConfig(),
		MinIO:       *LoadMinIOConfig(),
	}
	return cfg
}

func LoadJWTConfig() *JWTConfig {
	return &JWTConfig{
		Secret: shared.SafeGetEnv("JWT_SECRET", "your-secret-key"),
	}
}

func LoadDBConfig() *DBConfig {
	return &DBConfig{
		User:     shared.SafeGetEnv("DB_USER", "postgres"),
		Password: shared.SafeGetEnv("DB_PASSWORD", "postgres"),
		Name:     shared.SafeGetEnv("DB_NAME", "belimang"),
		Host:     shared.SafeGetEnv("DB_HOST", "localhost"),
		Port:     shared.SafeGetEnv("DB_PORT", "5432"),
		SSLMode:  shared.SafeGetEnv("DB_SSLMODE", "disable"),
		URL:      shared.SafeGetEnv("DATABASE_URL", ""),
	}
}

func LoadMinIOConfig() *MinIOConfig {
	return &MinIOConfig{
		Endpoint:        shared.SafeGetEnv("MINIO_ENDPOINT", "localhost:9000"),
		AccessKeyID:     shared.SafeGetEnv("MINIO_ACCESS_KEY", "minioadmin"),
		SecretAccessKey: shared.SafeGetEnv("MINIO_SECRET_KEY", "minioadmin123"),
		UseSSL:          shared.SafeGetEnv("MINIO_USE_SSL", "false") == "true",
		BucketName:      shared.SafeGetEnv("MINIO_BUCKET_NAME", "belimang-files"),
	}
}
