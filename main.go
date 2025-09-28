package main

import (
	"context"
	"time"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/config"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/db"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
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

	router := setupGin(cfg, pool)
	router.Run(":" + cfg.Port)
}

func setupGin(cfg *config.Config, pool *pgxpool.Pool) *gin.Engine {
	router := gin.New()

	// TODO: Add recovery middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	// router.Use(middleware.Recovery())

	// TODO: Add Route handlers
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Database ping
	router.GET("/db/ping", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"ok": true})
	})

	// Simple user creation (debug)
	router.POST("/debug/users", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
			Email    string `json:"email" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		queries := db.New(pool)
		if err := queries.CreateUser(c.Request.Context(), db.CreateUserParams{
			Username: req.Username,
			Password: req.Password,
			Email:    req.Email,
		}); err != nil {
			c.JSON(500, gin.H{"ok": false, "error": err.Error()})
			return
		}

		c.JSON(201, gin.H{"ok": true})
	})

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

	pgxConfig.MaxConns = 10
	pgxConfig.MinConns = 2
	pgxConfig.MaxConnLifetime = 30 * time.Minute
	pgxConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		log.Fatal().Msgf("Failed to create database pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		log.Fatal().Msgf("Failed to ping database: %v", err)
	}

	return pool
}
