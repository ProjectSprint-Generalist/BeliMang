package handlers

import (
	"context"
	"net/http"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/db"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/shared"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MerchantHandler struct {
	pool *pgxpool.Pool
}

func NewMerchantHandler(pool *pgxpool.Pool) *MerchantHandler {
	return &MerchantHandler{pool: pool}
}

func (h *MerchantHandler) CreateMerchant(c *gin.Context) {
	var payload dto.MerchantCreateRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid input: please make sure you have provided a valid name, merchant category, image URL, and location",
			Code:    http.StatusBadRequest,
		})
		return
	}

	queries := db.New(h.pool)
	ctx := context.Background()

	id, err := queries.CreateMerchant(ctx, db.CreateMerchantParams{
		Name:             payload.Name,
		MerchantCategory: db.MerchantCategory(payload.MerchantCategory),
		ImageUrl:         payload.ImageURL,
		Location:         payload.Location,
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

	c.JSON(http.StatusCreated, dto.MerchantCreateResponse{
		MerchantId: id.String(),
	})

}
