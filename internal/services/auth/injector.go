package auth

import (
	"github.com/ProjectSprint-Generalist/BeliMang/config"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth/handler"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth/repository"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Routes struct {
	handler *handler.AuthHandler
}

func NewRoutes(pool *pgxpool.Pool, jwtConfig config.JWTConfig) *Routes {
	repo := repository.NewUserRepositoryImpl(pool)
	svc := service.NewAuthService(repo, jwtConfig)
	return &Routes{
		handler: handler.NewAuthHandler(svc),
	}
}

func (r *Routes) Register(router *gin.Engine) {
	SetupRoutes(router, r.handler)
}
