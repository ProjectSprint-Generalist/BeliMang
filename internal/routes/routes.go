package routes

import (
	"github.com/ProjectSprint-Generalist/BeliMang/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, adminHandler *handlers.AdminHandler, adminMerchantHandler *handlers.AdminMerchantHandler) {
	admin := router.Group("/admin")
	{
		admin.POST("/register", adminHandler.RegisterAdmin)
		// admin.POST("/login", ...)

		admin.GET("/merchants", adminMerchantHandler.GetMerchants)
	}

	// users := router.Group("/users")
	// {
	// 	users.POST("/register", ...)
	// 	users.POST("/login", ...)
	// 	...
	// }

	// router.POST("/image", ...)

	// router.GET("/merchants/nearby/:lat/:long", ...)
}
