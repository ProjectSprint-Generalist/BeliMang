package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
	admin := router.Group("/admin")
	{
		// admin.POST("/register", ...)
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
