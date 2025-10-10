package order

import (
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/order/handler"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/order/repository"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/order/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Routes contains the order service dependencies
type Routes struct {
	handler *handler.OrderHandler
}

// NewRoutes creates and initializes all order service dependencies
func NewRoutes(pool *pgxpool.Pool, jwtSecret string) *Routes {
	// Initialize repository
	repo := repository.NewOrderRepositoryImpl(pool)

	// Initialize service
	svc := service.NewOrderService(repo)

	// Initialize handler
	h := handler.NewOrderHandler(svc)

	return &Routes{
		handler: h,
	}
}

// Register registers order routes with the provided router
func (r *Routes) Register(router *gin.Engine, jwtSecret string) {
	SetupRoutes(router, r.handler, jwtSecret)
}
