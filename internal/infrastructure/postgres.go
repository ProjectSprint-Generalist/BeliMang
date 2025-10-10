package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Postgres struct {
	Pool    *pgxpool.Pool
	cleanUp func()
}

func NewPostgres(ctx context.Context, url string) (*Postgres, error) {
	const (
		MaxConns              = 200
		MinConns              = 50
		MaxConnIdleTime       = 10 * time.Minute
		MaxConnLifetime       = 30 * time.Minute
		MaxConnLifetimeJitter = 2 * time.Minute
		HealthCheckPeriod     = 30 * time.Second
		ConnectTimeout        = 5 * time.Second
	)

	pgxConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	pgxConfig.MaxConns = MaxConns
	pgxConfig.MinConns = MinConns
	pgxConfig.MaxConnIdleTime = MaxConnIdleTime
	pgxConfig.MaxConnLifetime = MaxConnLifetime
	pgxConfig.MaxConnLifetimeJitter = MaxConnLifetimeJitter
	pgxConfig.HealthCheckPeriod = HealthCheckPeriod
	pgxConfig.ConnConfig.ConnectTimeout = ConnectTimeout

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	pg := &Postgres{Pool: pool, cleanUp: func() { pool.Close() }}

	if err := pg.runMigrations(url); err != nil {
		pg.cleanUp()
		return nil, fmt.Errorf("migrations failed: %w", err)
	}

	log.Info().Msg("Postgres ready")
	return pg, nil
}

func (p *Postgres) runMigrations(url string) error {
	m, err := migrate.New("file://migrations/", url)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func (p *Postgres) Close() error {
	if p.cleanUp != nil {
		p.cleanUp()
	}
	return nil
}
