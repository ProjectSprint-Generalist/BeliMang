package repository

import (
	"context"
	"errors"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/db"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryImpl struct {
	db *pgxpool.Pool
	q  *db.Queries
}

func NewUserRepositoryImpl(pool *pgxpool.Pool) UserRepository {
	return &UserRepositoryImpl{
		db: pool,
		q:  db.New(pool),
	}
}

func (r *UserRepositoryImpl) Save(ctx context.Context, u domain.User) error {
	err := r.q.CreateUser(ctx, db.CreateUserParams{
		Username: u.Username,
		Password: string(u.PasswordHash),
		Email:    u.Email,
		Role:     db.UserRole(u.Role),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrUserAlreadyExists
			}
			return ErrInternalServerError
		}
	}
	return nil
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string, role domain.Role) (domain.User, error) {
	user, err := r.q.GetUserByUsername(ctx, db.GetUserByUsernameParams{
		Username: username,
		Role:     db.UserRole(role),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.User{}, ErrUserNotFound
			}
			return domain.User{}, ErrInternalServerError
		}
		return domain.User{}, err
	}
	return domain.User{
		ID:           user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: []byte(user.Password),
		Role:         domain.Role(user.Role),
	}, nil
}
