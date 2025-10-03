package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/db"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/shared"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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

func (h *MerchantHandler) GetMerchants(c *gin.Context) {
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
		statusCode, errorMessage := shared.ParseDBResult(err)
		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Error:   errorMessage,
			Code:    statusCode,
		})
		return
	}

	// Get merchants
	merchants, err := queries.GetMerchants(ctx, queryParams)
	if err != nil {
		statusCode, errorMessage := shared.ParseDBResult(err)
		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Error:   errorMessage,
			Code:    statusCode,
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
			Location: dto.Location{
				Lat:  m.Lat,
				Long: m.Long,
			},
			CreatedAt: m.CreatedAt.Time.Format(shared.ISO8601WithNanoseconds),
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

func (h *MerchantHandler) CreateMerchantItem(c *gin.Context) {
	merchantID := c.Param("merchantId")

	var payload dto.MerchantItemCreateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid input: please make sure you have provided valid name, product category, price, and image URL",
			Code:    http.StatusBadRequest,
		})
		return
	}

	queries := db.New(h.pool)
	ctx := context.Background()

	// Validate merchantId format
	var merchantUUID pgtype.UUID
	if err := merchantUUID.Scan(merchantID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid merchant ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Check if merchant exists
	exists, err := queries.GetMerchantByID(ctx, merchantID)
	if err != nil {
		statusCode, errorMessage := shared.ParseDBResult(err)
		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Error:   errorMessage,
			Code:    statusCode,
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "Merchant not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	// Create merchant item
	itemID, err := queries.CreateMerchantItem(ctx, db.CreateMerchantItemParams{
		MerchantID:      merchantUUID,
		Name:            payload.Name,
		ProductCategory: db.ProductCategory(payload.ProductCategory),
		Price:           int32(payload.Price),
		ImageUrl:        payload.ImageURL,
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

	c.JSON(http.StatusCreated, dto.MerchantItemCreateResponse{
		ItemId: itemID.String(),
	})
}

func (h *MerchantHandler) GetMerchantItems(c *gin.Context) {
	merchantID := c.Param("merchantId")

	// Parse query parameters
	itemID := c.Query("itemId")
	name := c.Query("name")
	productCategory := c.Query("productCategory")
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

	// Validate merchantId format
	var merchantUUID pgtype.UUID
	if err := merchantUUID.Scan(merchantID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid merchant ID format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Check if merchant exists
	exists, err := queries.GetMerchantByID(ctx, merchantID)
	if err != nil {
		statusCode, errorMessage := shared.ParseDBResult(err)
		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Error:   errorMessage,
			Code:    statusCode,
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "Merchant not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	// Validate productCategory if provided
	validCategories := map[string]bool{
		"Beverage":  true,
		"Food":      true,
		"Snack":     true,
		"Condiment": true,
		"Additions": true,
	}

	if productCategory != "" && !validCategories[productCategory] {
		c.JSON(http.StatusOK, dto.GetMerchantItemsResponse{
			Data: []dto.MerchantItemData{},
			Meta: dto.MerchantMeta{
				Limit:  int(limit),
				Offset: int(offset),
				Total:  0,
			},
		})
		return
	}

	// Validate itemId if provided
	if itemID != "" {
		var tempUUID pgtype.UUID
		if err := tempUUID.Scan(itemID); err != nil {
			c.JSON(http.StatusOK, dto.GetMerchantItemsResponse{
				Data: []dto.MerchantItemData{},
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
	queryParams := db.MerchantItemQueryParams{
		MerchantID:      merchantID,
		ItemID:          itemID,
		Name:            name,
		ProductCategory: productCategory,
		CreatedAt:       createdAt,
		Limit:           limit,
		Offset:          offset,
	}

	// Get total count
	total, err := queries.CountMerchantItems(ctx, queryParams)
	if err != nil {
		statusCode, errorMessage := shared.ParseDBResult(err)
		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Error:   errorMessage,
			Code:    statusCode,
		})
		return
	}

	// Get merchant items
	items, err := queries.GetMerchantItems(ctx, queryParams)
	if err != nil {
		statusCode, errorMessage := shared.ParseDBResult(err)
		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Error:   errorMessage,
			Code:    statusCode,
		})
		return
	}

	// Convert to response format
	itemData := make([]dto.MerchantItemData, 0, len(items))
	for _, item := range items {
		// Convert UUID bytes to string
		uuidStr := ""
		if item.ID.Valid {
			uuidStr = pgtype.UUID{Bytes: item.ID.Bytes, Valid: true}.String()
		}

		itemData = append(itemData, dto.MerchantItemData{
			ItemId:          uuidStr,
			Name:            item.Name,
			ProductCategory: string(item.ProductCategory),
			Price:           int(item.Price),
			ImageURL:        item.ImageURL,
			CreatedAt:       item.CreatedAt.Time.Format(shared.ISO8601WithNanoseconds),
		})
	}

	c.JSON(http.StatusOK, dto.GetMerchantItemsResponse{
		Data: itemData,
		Meta: dto.MerchantMeta{
			Limit:  int(limit),
			Offset: int(offset),
			Total:  total,
		},
	})
}
