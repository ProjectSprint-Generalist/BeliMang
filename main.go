package main

import (
	"context"
	"time"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/config"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/handlers"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/middleware"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/routes"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/storage"
	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Error().Msg("No .env file found, using system environment variables")
	}

	cfg := config.LoadConfig()

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	pool := setupDatabase(cfg)

	defer pool.Close()

	minioClient, err := storage.NewMinioClient(&cfg.MinIO)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init minio")
	}

	router := setupGin(cfg, pool, minioClient)
	router.Run(":" + cfg.Port)
}

func setupGin(cfg *config.Config, pool *pgxpool.Pool, minioClient *storage.MinioClient) *gin.Engine {
	router := gin.New()

	// TODO: Add recovery middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	// router.Use(middleware.Recovery())

	adminHandler := handlers.NewAdminHandler(pool)
	userHandler := handlers.NewUserHandler(pool)
	merchantHandler := handlers.NewMerchantHandler(pool)
	imageHandler := handlers.NewImageHandler(pool, minioClient)
	estimateHandler := handlers.NewEstimateHandler(pool)
	routes.SetupRoutes(router, adminHandler, userHandler, merchantHandler, imageHandler, estimateHandler)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	return router
}

func setupDatabase(cfg *config.Config) *pgxpool.Pool {
	ctx := context.Background()

	if cfg.DB.URL == "" {
		log.Fatal().Msg("Database URL is not set.")
	}

	pgxConfig, err := pgxpool.ParseConfig(cfg.DB.URL)
	if err != nil {
		log.Fatal().Msgf("Failed to parse database URL: %v", err)
	}

	pgxConfig.MaxConns = 50
	pgxConfig.MinConns = 20
	pgxConfig.MaxConnIdleTime = 5 * time.Minute
	pgxConfig.MaxConnLifetimeJitter = 1 * time.Minute
	pgxConfig.MaxConnLifetime = 10 * time.Minute
	pgxConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatal().Msgf("Failed to create database pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		log.Fatal().Msgf("Failed to ping database: %v", err)
	}

	migrateDatabase(cfg)

	return pool
}

func migrateDatabase(cfg *config.Config) {
	m, err := migrate.New(
		"file://migrations/",
		cfg.DB.URL,
	)
	if err != nil {
		log.Fatal().Msgf("Failed to create migrate instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msgf("Failed to run migrations: %v", err)
	}
}
