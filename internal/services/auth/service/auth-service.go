package service

import (
	"context"
	"errors"
	"time"

	"github.com/ProjectSprint-Generalist/BeliMang/config"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth/domain"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth/repository"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateUser      = errors.New("user already exists")
)

type AuthService struct {
	repo   repository.UserRepository
	config config.JWTConfig
}

func NewAuthService(repo repository.UserRepository, cfg config.JWTConfig) *AuthService {
	return &AuthService{repo: repo, config: cfg}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string, role domain.Role) (string, error) {
	u, err := domain.NewUser(username, email, password, role)
	if err != nil {
		return "", err
	}
	if err := s.repo.Save(ctx, u); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return "", ErrDuplicateUser
		}
		return "", err
	}
	return s.GenerateToken(u)
}

func (s *AuthService) Login(ctx context.Context, username, password string, role domain.Role) (string, error) {
	u, err := s.repo.GetByUsername(ctx, username, role)
	if err != nil {
		return "", err
	}
	if !u.VerifyPassword(password) {
		return "", ErrInvalidCredentials
	}
	return s.GenerateToken(u)
}

func (s *AuthService) GenerateToken(u domain.User) (string, error) {
	claims := &domain.JWTClaim{
		UserID:   u.ID,
		Username: u.Username,
		Email:    u.Email,
		Role:     string(u.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Secret))
}
