package routes

import (
	"github.com/ProjectSprint-Generalist/BeliMang/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, adminHandler *handlers.AdminHandler) {
	admin := router.Group("/admin")
	{
		admin.POST("/register", adminHandler.RegisterAdmin)
		// admin.POST("/login", ...)
		// ...
	}

	users := router.Group("/users")
	{
		// users.POST("/register", ...)
		// users.POST("/login", ...)
		// ...
	}

	// router.POST("/image", ...)

	// router.GET("/merchants/nearby/:lat/:long", ...)
}
