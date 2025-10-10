package auth

import (
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authHandler *handler.AuthHandler) {
	admin := router.Group("/admin")
	{
		admin.POST("/register", authHandler.RegisterAdmin)
		admin.POST("/login", authHandler.LoginAdmin)
	}

	user := router.Group("/users")
	{
		user.POST("/register", authHandler.RegisterUser)
		user.POST("/login", authHandler.LoginUser)
	}
}
