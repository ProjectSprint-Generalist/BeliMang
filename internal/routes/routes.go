package routes

import (
	"github.com/ProjectSprint-Generalist/BeliMang/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	adminHandler *handlers.AdminHandler,
	userHandler *handlers.UserHandler,
	merchantHandler *handlers.MerchantHandler,
) {
	admin := router.Group("/admin")
	{
		admin.POST("/register", adminHandler.RegisterAdmin)
		// admin.POST("/login", ...)
		// ...

		merchant := admin.Group("/merchants")
		{
			merchant.POST("/", merchantHandler.CreateMerchant)
		}
	}

	users := router.Group("/users")
	{
		users.POST("/register", userHandler.RegisterUser)
		users.POST("/login", userHandler.LoginUser)
	}

	// router.POST("/image", ...)

	// router.GET("/merchants/nearby/:lat/:long", ...)
}
