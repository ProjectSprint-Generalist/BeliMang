package routes

import (
	"github.com/ProjectSprint-Generalist/BeliMang/internal/handlers"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, adminHandler *handlers.AdminHandler, userHandler *handlers.UserHandler, merchantHandler *handlers.MerchantHandler, imageHandler *handlers.ImageHandler) {
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

	image := router.Group("/image")
	image.Use(middleware.AuthMiddleware(), middleware.IsAuthorized("admin"))
	{
		image.POST("", imageHandler.UploadImage)
	}

	// router.GET("/merchants/nearby/:lat/:long", ...)
}
