package repository

import (
	"context"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth/domain"
)

type UserRepository interface {
	Save(ctx context.Context, u domain.User) error
	GetByUsername(ctx context.Context, username string, role domain.Role) (domain.User, error)
}
