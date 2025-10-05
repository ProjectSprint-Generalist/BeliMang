package routes

import (
	"github.com/ProjectSprint-Generalist/BeliMang/internal/handlers"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, adminHandler *handlers.AdminHandler, userHandler *handlers.UserHandler, merchantHandler *handlers.MerchantHandler, imageHandler *handlers.ImageHandler, estimateHandler *handlers.EstimateHandler, orderHandler *handlers.OrderHandler) {
	admin := router.Group("/admin")
	{
		admin.POST("/register", adminHandler.RegisterAdmin)
		admin.POST("/login", adminHandler.LoginAdmin)

		merchant := admin.Group("/merchants")
		merchant.Use(middleware.AuthMiddleware(), middleware.IsAuthorized("admin"))
		{
			merchant.GET("", merchantHandler.GetMerchants)
			merchant.POST("/", merchantHandler.CreateMerchant)
			merchant.GET("/:merchantId/items", merchantHandler.GetMerchantItems)
			merchant.POST("/:merchantId/items", merchantHandler.CreateMerchantItem)
		}
	}

	users := router.Group("/users")
	{
		users.POST("/register", userHandler.RegisterUser)
		users.POST("/login", userHandler.LoginUser)
		users.POST("/estimate", middleware.AuthMiddleware(), middleware.IsAuthorized("user"), estimateHandler.Estimate)
		users.POST("/orders", middleware.AuthMiddleware(), middleware.IsAuthorized("user"), orderHandler.CreateOrder)
		users.GET("/orders", middleware.AuthMiddleware(), middleware.IsAuthorized("user"), orderHandler.GetOrders)
	}

	image := router.Group("/image")
	image.Use(middleware.AuthMiddleware(), middleware.IsAuthorized("admin"))
	{
		image.POST("", imageHandler.UploadImage)
	}

	// Nearby merchants endpoint
	merchants := router.Group("/merchants")
	merchants.Use(middleware.AuthMiddleware(), middleware.IsAuthorized("user"))
	{
		// Path pattern: /merchants/nearby/:coords where :coords is "lat,long"
		merchants.GET("/nearby/:coords", merchantHandler.GetNearbyMerchants)
	}
}
