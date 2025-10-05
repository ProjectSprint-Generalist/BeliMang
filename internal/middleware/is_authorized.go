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

// GetAuthUser retrieves the authenticated user from the Gin context
func GetAuthUser(c *gin.Context) (dto.AuthUser, bool) {
	username, userExists := c.Get("username")
	role, roleExists := c.Get("role")
	email, emailExists := c.Get("email")

	if !userExists || !roleExists || !emailExists {
		return dto.AuthUser{}, false
	}

	return dto.AuthUser{
		Username: username.(string),
		Role:     role.(string),
		Email:    email.(string),
	}, true
}
