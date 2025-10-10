package merchant

import (
	"github.com/ProjectSprint-Generalist/BeliMang/internal/server/middleware"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth/domain"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/merchant/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, merchantHandler *handler.MerchantHandler, jwtSecret string) {
	admin := router.Group("/admin")
	admin.Use(middleware.RequireAuth(jwtSecret), middleware.IsAuthorized(domain.RoleAdmin))
	{
		merchant := admin.Group("/merchants")
		{
			merchant.GET("", merchantHandler.GetMerchants)
			merchant.POST("/", merchantHandler.CreateMerchant)
			merchant.GET("/:merchantId/items", merchantHandler.GetMerchantItems)
			merchant.POST("/:merchantId/items", merchantHandler.CreateMerchantItem)
		}
	}

	merchants := router.Group("/merchants")
	merchants.Use(middleware.RequireAuth(jwtSecret), middleware.IsAuthorized(domain.RoleUser))
	{
		merchants.GET("/nearby/:coords", merchantHandler.GetNearbyMerchants)
	}
}
