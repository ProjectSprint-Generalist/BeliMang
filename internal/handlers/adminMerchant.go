package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/db"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AdminMerchantHandler wires admin endpoints to sqlc-generated queries.
type AdminMerchantHandler struct {
	pool *pgxpool.Pool
}

func NewAdminMerchantHandler(pool *pgxpool.Pool) *AdminMerchantHandler {
	return &AdminMerchantHandler{pool: pool}
}

func (h *AdminMerchantHandler) GetMerchants(c *gin.Context) {
	// Parse query parameters
	merchantID := c.Query("merchantId")
	name := c.Query("name")
	merchantCategory := c.Query("merchantCategory")
	createdAt := c.Query("createdAt")

	// Parse limit and offset with defaults
	limit := int32(5)
	if limitStr := c.Query("limit"); limitStr != "" {
		if val, err := strconv.ParseInt(limitStr, 10, 32); err == nil && val > 0 {
			limit = int32(val)
		}
	}

	offset := int32(0)
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if val, err := strconv.ParseInt(offsetStr, 10, 32); err == nil && val >= 0 {
			offset = int32(val)
		}
	}

	// Default sort order is desc
	if createdAt != "asc" && createdAt != "desc" {
		createdAt = "desc"
	}

	queries := db.New(h.pool)
	ctx := context.Background()

	// Validate merchantCategory if provided
	validCategories := map[string]bool{
		"SmallRestaurant":       true,
		"MediumRestaurant":      true,
		"LargeRestaurant":       true,
		"MerchandiseRestaurant": true,
		"BoothKiosk":            true,
		"ConvenienceStore":      true,
	}

	if merchantCategory != "" && !validCategories[merchantCategory] {
		c.JSON(http.StatusOK, dto.GetMerchantsResponse{
			Data: []dto.MerchantData{},
			Meta: dto.MerchantMeta{
				Limit:  int(limit),
				Offset: int(offset),
				Total:  0,
			},
		})
		return
	}

	// Validate merchantId if provided
	if merchantID != "" {
		var tempUUID pgtype.UUID
		if err := tempUUID.Scan(merchantID); err != nil {
			c.JSON(http.StatusOK, dto.GetMerchantsResponse{
				Data: []dto.MerchantData{},
				Meta: dto.MerchantMeta{
					Limit:  int(limit),
					Offset: int(offset),
					Total:  0,
				},
			})
			return
		}
	}

	// Build query params
	queryParams := db.MerchantQueryParams{
		MerchantID:       merchantID,
		Name:             name,
		MerchantCategory: merchantCategory,
		CreatedAt:        createdAt,
		Limit:            limit,
		Offset:           offset,
	}

	// Get total count
	total, err := queries.CountMerchants(ctx, queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to fetch merchants count",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Get merchants
	merchants, err := queries.GetMerchants(ctx, queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to fetch merchants",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Convert to response format
	merchantData := make([]dto.MerchantData, 0, len(merchants))
	for _, m := range merchants {
		// Convert UUID bytes to string
		uuidStr := ""
		if m.ID.Valid {
			uuidStr = pgtype.UUID{Bytes: m.ID.Bytes, Valid: true}.String()
		}

		merchantData = append(merchantData, dto.MerchantData{
			MerchantID:       uuidStr,
			Name:             m.Name,
			MerchantCategory: string(m.MerchantCategory),
			ImageURL:         m.ImageURL,
			Location: dto.MerchantLocation{
				Lat:  m.Lat,
				Long: m.Long,
			},
			CreatedAt: m.CreatedAt.Time.Format("2006-01-02T15:04:05.999999999Z07:00"),
		})
	}

	c.JSON(http.StatusOK, dto.GetMerchantsResponse{
		Data: merchantData,
		Meta: dto.MerchantMeta{
			Limit:  int(limit),
			Offset: int(offset),
			Total:  total,
		},
	})
}
