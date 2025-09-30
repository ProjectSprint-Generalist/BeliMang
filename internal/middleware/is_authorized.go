package middleware

import (
	"net/http"
	"strings"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/gin-gonic/gin"
)

func IsAuthorized() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" {

			response := dto.ErrorResponse{
				Success: false,
				Error:   "Authorization header required",
				Code:    http.StatusUnauthorized,
			}
			context.JSON(http.StatusUnauthorized, response)
			context.Abort()
			return
		}

		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			response := dto.ErrorResponse{
				Success: false,
				Error:   "Invalid or expired token",
				Code:    http.StatusUnauthorized,
			}
			context.JSON(http.StatusUnauthorized, response)
			context.Abort()
			return
		}

		context.Set("username", claims.Username)
		context.Set("role", claims.Role)

		context.Next()
	}
}
