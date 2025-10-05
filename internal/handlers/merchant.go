package handlers

import (
	"context"
	"net/http"
	"net/url"
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

// isValidImageURL validates if a URL is a proper image URL
func isValidImageURL(imageURL string) bool {
	if imageURL == "" {
		return false
	}

	// Parse the URL
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return false
	}

	// Check if scheme is http or https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	// Check if host is present
	if parsedURL.Host == "" {
		return false
	}

	// Check if path is present (for images, we want at least a path)
	if parsedURL.Path == "" || parsedURL.Path == "/" {
		return false
	}

	return true
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

	// Validate merchant category manually since it's a custom type
	if !dto.ValidMerchantCategories[payload.MerchantCategory] {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid merchant category. Must be one of: SmallRestaurant, MediumRestaurant, LargeRestaurant, MerchandiseRestaurant, BoothKiosk, ConvenienceStore",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate location coordinates manually (nested struct validation may not work reliably)
	if payload.Location.Lat < -90 || payload.Location.Lat > 90 || payload.Location.Lat == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid latitude. Must be between -90 and 90, and not zero",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if payload.Location.Long < -180 || payload.Location.Long > 180 || payload.Location.Long == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid longitude. Must be between -180 and 180, and not zero",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate image URL manually (Gin's url validator is too permissive)
	if !isValidImageURL(payload.ImageURL) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid image URL. Must be a complete HTTP/HTTPS URL with a path (e.g., https://example.com/image.jpg)",
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
		StMakepoint:      payload.Location.Long,
		StMakepoint_2:    payload.Location.Lat,
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

	// Build query params for generated functions
	var merchantIDText pgtype.Text
	if merchantID != "" {
		merchantIDText = pgtype.Text{String: merchantID, Valid: true}
	}
	var nameText pgtype.Text
	if name != "" {
		nameText = pgtype.Text{String: name, Valid: true}
	}
	var categoryText pgtype.Text
	if merchantCategory != "" {
		categoryText = pgtype.Text{String: merchantCategory, Valid: true}
	}

	// Get total count
	total, err := queries.CountMerchants(ctx, db.CountMerchantsParams{
		MerchantID:       merchantIDText,
		MerchantCategory: categoryText,
		Name:             nameText,
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

	// Get merchants
	merchants, err := queries.GetMerchants(ctx, db.GetMerchantsParams{
		MerchantID:       merchantIDText,
		MerchantCategory: categoryText,
		Name:             nameText,
		CreatedAt:        createdAt,
		OffsetVal:        offset,
		LimitVal:         limit,
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

	// Convert to response format
	merchantData := make([]dto.MerchantData, 0, len(merchants))
	for _, m := range merchants {
		// Convert UUID bytes to string
		uuidStr := ""
		if m.ID.Valid {
			uuidStr = pgtype.UUID{Bytes: m.ID.Bytes, Valid: true}.String()
		}

		// Handle lat/long conversion from interface{} to float64
		var lat, long float64
		if m.Lat != nil {
			if latVal, ok := m.Lat.(float64); ok {
				lat = latVal
			}
		}
		if m.Long != nil {
			if longVal, ok := m.Long.(float64); ok {
				long = longVal
			}
		}

		merchantData = append(merchantData, dto.MerchantData{
			MerchantID:       uuidStr,
			Name:             m.Name,
			MerchantCategory: string(m.MerchantCategory),
			ImageURL:         m.ImageUrl,
			Location: dto.Location{
				Lat:  lat,
				Long: long,
			},
			CreatedAt: m.CreatedAt.Time.Format(shared.ISO8601WithNanoseconds),
		})
	}

	c.JSON(http.StatusOK, dto.GetMerchantsResponse{
		Data: merchantData,
		Meta: dto.MerchantMeta{
			Limit:  int(limit),
			Offset: int(offset),
			Total:  int(total),
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
	exists, err := queries.GetMerchantByID(ctx, merchantUUID)
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
	exists, err := queries.GetMerchantByID(ctx, merchantUUID)
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
		"Condiments": true,
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

	// Build query params for generated functions
	var itemIDText pgtype.Text
	if itemID != "" {
		itemIDText = pgtype.Text{String: itemID, Valid: true}
	}
	var nameText pgtype.Text
	if name != "" {
		nameText = pgtype.Text{String: name, Valid: true}
	}
	var productCategoryText pgtype.Text
	if productCategory != "" {
		productCategoryText = pgtype.Text{String: productCategory, Valid: true}
	}

	// Get total count
	total, err := queries.CountMerchantItems(ctx, db.CountMerchantItemsParams{
		MerchantID:      merchantUUID,
		ItemID:          itemIDText,
		ProductCategory: productCategoryText,
		Name:            nameText,
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

	// Get merchant items
	items, err := queries.GetMerchantItems(ctx, db.GetMerchantItemsParams{
		MerchantID:      merchantUUID,
		ItemID:          itemIDText,
		ProductCategory: productCategoryText,
		Name:            nameText,
		CreatedAt:       createdAt,
		OffsetVal:       offset,
		LimitVal:        limit,
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
			ImageURL:        item.ImageUrl,
			CreatedAt:       item.CreatedAt.Time.Format(shared.ISO8601WithNanoseconds),
		})
	}

	c.JSON(http.StatusOK, dto.GetMerchantItemsResponse{
		Data: itemData,
		Meta: dto.MerchantMeta{
			Limit:  int(limit),
			Offset: int(offset),
			Total:  int(total),
		},
	})
}

func (h *MerchantHandler) GetNearbyMerchants(c *gin.Context) {
	// Path param in the form ":coords" where value is "lat,long"
	coords := c.Param("coords")
	var lat, long float64
	if coords != "" {
		// Expecting "lat,long"
		commaIdx := -1
		for i := 0; i < len(coords); i++ {
			if coords[i] == ',' {
				commaIdx = i
				break
			}
		}
		if commaIdx <= 0 || commaIdx >= len(coords)-1 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Error:   "lat/long is not valid",
				Code:    http.StatusBadRequest,
			})
			return
		}
		latStr := coords[:commaIdx]
		longStr := coords[commaIdx+1:]
		latParsed, err1 := strconv.ParseFloat(latStr, 64)
		longParsed, err2 := strconv.ParseFloat(longStr, 64)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Error:   "lat/long is not valid",
				Code:    http.StatusBadRequest,
			})
			return
		}
		if latParsed < -90 || latParsed > 90 || longParsed < -180 || longParsed > 180 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Error:   "lat/long is not valid",
				Code:    http.StatusBadRequest,
			})
			return
		}
		lat = latParsed
		long = longParsed
	} else {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "lat/long is not valid",
			Code:    http.StatusBadRequest,
		})
		return
	}

	merchantId := c.Query("merchantId")
	name := c.Query("name")
	merchantCategory := c.Query("merchantCategory")

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
		c.JSON(http.StatusOK, dto.GetNearbyMerchantsResponse{
			Data: []dto.NearbyMerchant{},
			Meta: dto.MerchantMeta{
				Limit:  int(limit),
				Offset: int(offset),
				Total:  0,
			},
		})
		return
	}

	// Validate merchantId if provided (invalid -> 200 empty)
	if merchantId != "" {
		var tempUUID pgtype.UUID
		if err := tempUUID.Scan(merchantId); err != nil {
			c.JSON(http.StatusOK, dto.GetNearbyMerchantsResponse{
				Data: []dto.NearbyMerchant{},
				Meta: dto.MerchantMeta{
					Limit:  int(limit),
					Offset: int(offset),
					Total:  0,
				},
			})
			return
		}
	}

	// Prepare optional filter params for sqlc narg
	var merchantIDText pgtype.Text
	if merchantId != "" {
		merchantIDText = pgtype.Text{String: merchantId, Valid: true}
	}
	var nameText pgtype.Text
	if name != "" {
		nameText = pgtype.Text{String: name, Valid: true}
	}
	var categoryText pgtype.Text
	if merchantCategory != "" {
		categoryText = pgtype.Text{String: merchantCategory, Valid: true}
	}

	total, err := queries.CountNearbyMerchants(ctx, db.CountNearbyMerchantsParams{
		MerchantID:       merchantIDText,
		MerchantCategory: categoryText,
		Name:             nameText,
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

	rows, err := queries.GetNearbyMerchants(ctx, db.GetNearbyMerchantsParams{
		Lat:              lat,
		Long:             long,
		MerchantID:       merchantIDText,
		MerchantCategory: categoryText,
		Name:             nameText,
		RowLimit:         limit,
		RowOffset:        offset,
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

	// Assemble response with items per merchant
	resp := make([]dto.NearbyMerchant, 0, len(rows))
	for _, m := range rows {
		// Convert UUID to string
		merchantIDStr := ""
		if m.ID.Valid {
			merchantIDStr = pgtype.UUID{Bytes: m.ID.Bytes, Valid: true}.String()
		}

		lat64, _ := m.Lat.(float64)
		long64, _ := m.Long.(float64)

		merchantData := dto.MerchantData{
			MerchantID:       merchantIDStr,
			Name:             m.Name,
			MerchantCategory: string(m.MerchantCategory),
			ImageURL:         m.ImageUrl,
			Location: dto.Location{
				Lat:  lat64,
				Long: long64,
			},
			CreatedAt: m.CreatedAt.Time.Format(shared.ISO8601WithNanoseconds),
		}

		// Fetch items for this merchant (no extra filters), unpaged
		var merchantUUIDForItems pgtype.UUID
		if err := merchantUUIDForItems.Scan(merchantIDStr); err != nil {
			continue // Skip this merchant if UUID is invalid
		}

		items, err := queries.GetMerchantItems(ctx, db.GetMerchantItemsParams{
			MerchantID:      merchantUUIDForItems,
			ItemID:          pgtype.Text{}, // empty text for no filter
			ProductCategory: pgtype.Text{}, // empty text for no filter
			Name:            pgtype.Text{}, // empty text for no filter
			CreatedAt:       "desc",
			OffsetVal:       0,
			LimitVal:        1000,
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

		itemData := make([]dto.MerchantItemData, 0, len(items))
		for _, it := range items {
			itemIDStr := ""
			if it.ID.Valid {
				itemIDStr = pgtype.UUID{Bytes: it.ID.Bytes, Valid: true}.String()
			}
			itemData = append(itemData, dto.MerchantItemData{
				ItemId:          itemIDStr,
				Name:            it.Name,
				ProductCategory: string(it.ProductCategory),
				Price:           int(it.Price),
				ImageURL:        it.ImageUrl,
				CreatedAt:       it.CreatedAt.Time.Format(shared.ISO8601WithNanoseconds),
			})
		}

		resp = append(resp, dto.NearbyMerchant{
			Merchant: merchantData,
			Items:    itemData,
		})
	}

	c.JSON(http.StatusOK, dto.GetNearbyMerchantsResponse{
		Data: resp,
		Meta: dto.MerchantMeta{
			Limit:  int(limit),
			Offset: int(offset),
			Total:  int(total),
		},
	})
}
