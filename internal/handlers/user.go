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

// UserHandler wires user endpoints to sqlc-generated queries.
type UserHandler struct {
	pool *pgxpool.Pool
}

func NewUserHandler(pool *pgxpool.Pool) *UserHandler {
	return &UserHandler{pool: pool}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var payload dto.UserRegisterRequest

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
	err = queries.CreateUser(ctx, db.CreateUserParams{
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
		Role:     "user",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to generate authentication token",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, dto.UserAuthResponse{
		Token: token,
	})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var payload dto.UserLoginRequest

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

	fetchedUser, err := queries.GetUserByUsername(ctx, payload.Username)
	if err != nil {
		// Return 400 for invalid username/password as per requirements
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid username or password",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := shared.VerifyPassword(payload.Password, fetchedUser.Password); err != nil {
		// Return 400 for invalid password as per requirements
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid username or password",
			Code:    http.StatusBadRequest,
		})
		return
	}

	token, err := middleware.GenerateToken(dto.AuthUser{
		Username: fetchedUser.Username,
		Email:    fetchedUser.Email,
		Role:     "user",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to generate token",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, dto.UserAuthResponse{
		Token: token,
	})
}
