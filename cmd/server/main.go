package main

import (
	"context"
	"os"

	"github.com/ProjectSprint-Generalist/BeliMang/config"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/infrastructure"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/server"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using system environment variables")
	}

	cfg := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	deps, cleanup, err := infrastructure.Init(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize dependencies")
		os.Exit(1)
	}
	defer cleanup()

	server := server.NewServer(cfg, deps)
	server.Run()
}
