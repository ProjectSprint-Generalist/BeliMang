package handlers

import (
	"context"
	"net/http"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/db"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/middleware"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/shared"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AdminHandler wires admin endpoints to sqlc-generated queries.
type AdminHandler struct {
	pool *pgxpool.Pool
}

func NewAdminHandler(pool *pgxpool.Pool) *AdminHandler {
	return &AdminHandler{pool: pool}
}

func (h *AdminHandler) RegisterAdmin(c *gin.Context) {
	var payload dto.AdminRegisterRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid input: please make sure you have provided a valid username, email, and password",
			Code:    http.StatusBadRequest,
		})
		return
	}

	queries := db.New(h.pool)
	ctx := context.Background()

	hashedPassword, err := shared.HashPassword(payload.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to hash password",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Try to insert to db
	err = queries.CreateAdmin(ctx, db.CreateAdminParams{
		Username: payload.Username,
		Password: hashedPassword,
		Email:    payload.Email,
	})

	if err != nil {
		statusCode, errorMessage := shared.ParseDBResult(err)
		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Error:   errorMessage,
			Code:    statusCode,
		})
		return
	}

	// Generate JWT
	token, err := middleware.GenerateToken(dto.AuthUser{
		Username: payload.Username,
		Email:    payload.Email,
		Role:     "admin",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to generate authentication token",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, dto.AdminRegisterResponse{
		Token: token,
	})

}

func (h *AdminHandler) LoginAdmin(c *gin.Context) {
	var payload dto.AdminLoginRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid input: please make sure you have provided a valid username and password",
			Code:    http.StatusBadRequest,
		})
		return
	}

	queries := db.New(h.pool)
	ctx := context.Background()

	fetchedAdmin, err := queries.GetAdminByUsername(ctx, payload.Username)
	if err != nil {
		statusCode, errorMessage := shared.ParseDBResult(err)
		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Error:   errorMessage,
			Code:    statusCode,
		})
		return
	}

	if err := shared.VerifyPassword(payload.Password, fetchedAdmin.Password); err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid password",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	token, err := middleware.GenerateToken(dto.AuthUser{
		Username: fetchedAdmin.Password,
		Email:    fetchedAdmin.Email,
		Role:     "admin",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to generate token",
			Code:    http.StatusInternalServerError,
		})
	}

	c.JSON(http.StatusOK, dto.AdminLoginResponse{
		Token: token,
	})

}
