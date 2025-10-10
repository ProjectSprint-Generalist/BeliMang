package server

import (
	"github.com/ProjectSprint-Generalist/BeliMang/config"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/infrastructure"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/server/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
}

func NewServer(cfg *config.Config, deps *infrastructure.Dependencies) *Server {
	gin.SetMode(ginMode(cfg.Environment))

	r := gin.New()
	r.Use(middleware.Recovery(), middleware.Logger(), middleware.CORS())

	h := BuildHandlers(deps)
	h.Auth.Register(r)
	h.Merchant.Register(r, cfg.JWT.Secret)
	h.Order.Register(r, cfg.JWT.Secret)

	return &Server{router: r, cfg: cfg}
}

func (s *Server) Run() error {
	return s.router.Run(":" + s.cfg.Port)
}

func ginMode(env string) string {
	if env == "production" {
		return gin.ReleaseMode
	}
	return gin.DebugMode
}
