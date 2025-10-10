package middleware

import (
	"net/http"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Success: false,
				Error:   "Internal server error: " + err,
				Code:    http.StatusInternalServerError,
			})
		}
		c.Abort()
	})
}
