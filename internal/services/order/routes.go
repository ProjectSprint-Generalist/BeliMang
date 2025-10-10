package order

import (
	"github.com/ProjectSprint-Generalist/BeliMang/internal/server/middleware"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth/domain"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/order/handler"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures order-related routes
func SetupRoutes(router *gin.Engine, orderHandler *handler.OrderHandler, jwtSecret string) {
	users := router.Group("/users")
	users.Use(middleware.RequireAuth(jwtSecret), middleware.IsAuthorized(domain.RoleUser))
	{
		// Estimate endpoint
		users.POST("/estimate", orderHandler.Estimate)

		// Order endpoints
		users.POST("/orders", orderHandler.CreateOrder)
		users.GET("/orders", orderHandler.GetOrders)
	}
}
