package middleware

import (
	"net/http"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/gin-gonic/gin"
)

func IsAuthorized(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists || userRole != role {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Success: false,
				Error:   "Forbidden: insufficient role",
				Code:    http.StatusForbidden,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
