package infrastructure

import (
	"context"

	"github.com/ProjectSprint-Generalist/BeliMang/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	DBPool      *pgxpool.Pool
	MinioClient *MinioClient
	Cleanup     func()
}

func Init(ctx context.Context, cfg *config.Config) (*Dependencies, func(), error) {
	pg, err := NewPostgres(ctx, cfg.DB.URL)
	if err != nil {
		return nil, nil, err
	}

	minio, err := NewMinioClient(&cfg.MinIO)
	if err != nil {
		pg.Close()
		return nil, nil, err
	}

	deps := &Dependencies{
		DBPool:      pg.Pool,
		MinioClient: minio,
	}
	return deps, func() { pg.Close() }, nil
}
