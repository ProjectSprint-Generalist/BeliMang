package server

import (
	"github.com/ProjectSprint-Generalist/BeliMang/config"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/infrastructure"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/merchant"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/order"
)

type Handlers struct {
	Auth     *auth.Routes
	Merchant *merchant.Routes
	Order    *order.Routes
}

func BuildHandlers(deps *infrastructure.Dependencies) *Handlers {
	cfg := config.LoadConfig()
	return &Handlers{
		Auth:     auth.NewRoutes(deps.DBPool, cfg.JWT),
		Merchant: merchant.NewRoutes(deps.DBPool, deps.MinioClient, cfg.JWT.Secret),
		Order:    order.NewRoutes(deps.DBPool, cfg.JWT.Secret),
	}
}
