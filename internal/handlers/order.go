package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/db"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderHandler struct {
	Q *db.Queries
}

func NewOrderHandler(pool *pgxpool.Pool) *OrderHandler {
	q := db.New(pool)
	return &OrderHandler{Q: q}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Get user from JWT
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Get user ID from database
	user, err := h.Q.GetUserByUsername(c, username.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Validate that the calculated estimate exists and belongs to the user
	estimateUUID, err := uuid.Parse(req.CalculatedEstimateID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid calculatedEstimateId",
			Code:    http.StatusBadRequest,
		})
		return
	}

	estimateID := pgtype.UUID{Bytes: estimateUUID, Valid: true}
	estimate, err := h.Q.GetCalculatedEstimateByID(c, estimateID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   "Calculated estimate not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	// Create the order
	orderID, err := h.Q.CreateOrder(c, db.CreateOrderParams{
		UserID:               user.ID,
		CalculatedEstimateID: estimate.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to create order",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateOrderResponse{
		OrderID: orderID.String(),
	})
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	// Get user from JWT
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Get user ID from database
	user, err := h.Q.GetUserByUsername(c, username.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Parse query parameters
	var params dto.GetOrdersParams
	if err := c.ShouldBindQuery(&params); err != nil {
		params = dto.GetOrdersParams{
			Limit:  5,
			Offset: 0,
		}
	}

	// Set defaults
	if params.Limit <= 0 {
		params.Limit = 5
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	// Get total count for pagination metadata
	totalCount, err := h.Q.GetOrdersCountByUserID(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to get orders count",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Get orders from database
	orders, err := h.Q.GetOrdersByUserID(c, db.GetOrdersByUserIDParams{
		Column1: user.ID,
		Limit:   int32(params.Limit * 2), // Get more to account for filtering
		Offset:  int32(params.Offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to get orders",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Process the orders and build response
	response := h.buildOrdersResponse(c, orders, params, int(totalCount))

	c.JSON(http.StatusOK, response.Data)
}

func (h *OrderHandler) buildOrdersResponse(c *gin.Context, orders []db.GetOrdersByUserIDRow, params dto.GetOrdersParams, totalCount int) dto.GetOrdersResponse {
	var response dto.GetOrdersResponse
	response.Meta.Limit = params.Limit
	response.Meta.Offset = params.Offset

	// Will be updated based on filtered count
	response.Meta.Total = totalCount

	filteredCount := 0
	var allFilteredOrders []dto.OrderHistory // To track all filtered orders for accurate total
	var filteredData []dto.OrderHistory      // To store paginated results

	// First pass: Get all orders that match the filters
	for _, order := range orders {
		orderIDStr := order.ID.String()

		// Parse estimate data
		var estimateRequest dto.EstimateRequest
		if err := json.Unmarshal(order.EstimateData, &estimateRequest); err != nil {
			continue
		}

		// Apply filters
		if !h.matchesFilters(c, estimateRequest, params) {
			continue
		}

		// Extract order details from estimate data
		orderDetails, err := h.extractOrderDetails(c, estimateRequest)
		if err != nil {
			continue
		}

		orderHistory := dto.OrderHistory{
			OrderID: orderIDStr,
			Orders:  orderDetails,
		}

		allFilteredOrders = append(allFilteredOrders, orderHistory)
	}

	// Update total count based on filter results
	response.Meta.Total = len(allFilteredOrders)

	// Second pass: Apply pagination
	for i, orderHistory := range allFilteredOrders {
		// Skip orders before offset
		if i < params.Offset {
			continue
		}

		// Stop if we have enough results for this page
		if filteredCount >= params.Limit {
			break
		}

		filteredData = append(filteredData, orderHistory)
		filteredCount++
	}

	response.Data = filteredData
	return response
}

func (h *OrderHandler) matchesFilters(c *gin.Context, estimateRequest dto.EstimateRequest, params dto.GetOrdersParams) bool {
	// Apply merchant ID filter
	if params.MerchantID != nil && *params.MerchantID != "" {
		found := false
		for _, order := range estimateRequest.Orders {
			if order.MerchantId == *params.MerchantID {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Apply name filter
	if params.Name != nil && *params.Name != "" {
		// We need to check if merchant name or item name contains the search string
		nameMatch := false

		// Check each order for matching merchant or item name
		for _, order := range estimateRequest.Orders {
			// Get merchant details
			merchantUUID, err := uuid.Parse(order.MerchantId)
			if err != nil {
				continue
			}

			merchantPgUUID := pgtype.UUID{Bytes: merchantUUID, Valid: true}
			merchant, err := h.Q.GetMerchantDetailsByID(c, merchantPgUUID)
			if err != nil {
				continue
			}

			// Check merchant name (case insensitive)
			if strings.Contains(strings.ToLower(merchant.Name), strings.ToLower(*params.Name)) {
				nameMatch = true
				break
			}

			// Check item names
			for _, item := range order.Items {
				itemUUID, err := uuid.Parse(item.ItemId)
				if err != nil {
					continue
				}

				itemPgUUID := pgtype.UUID{Bytes: itemUUID, Valid: true}
				merchantItem, err := h.Q.GetMerchantItemByID(c, itemPgUUID)
				if err != nil {
					continue
				}

				// Check item name (case insensitive)
				if strings.Contains(strings.ToLower(merchantItem.Name), strings.ToLower(*params.Name)) {
					nameMatch = true
					break
				}
			}

			if nameMatch {
				break
			}
		}

		if !nameMatch {
			return false
		}
	}

	// Apply merchant category filter
	if params.MerchantCategory != nil && *params.MerchantCategory != "" {
		categoryMatch := false

		// Check each order for matching merchant category
		for _, order := range estimateRequest.Orders {
			// Get merchant details
			merchantUUID, err := uuid.Parse(order.MerchantId)
			if err != nil {
				continue
			}

			merchantPgUUID := pgtype.UUID{Bytes: merchantUUID, Valid: true}
			merchant, err := h.Q.GetMerchantDetailsByID(c, merchantPgUUID)
			if err != nil {
				continue
			}

			// Check merchant category
			if string(merchant.MerchantCategory) == *params.MerchantCategory {
				categoryMatch = true
				break
			}
		}

		if !categoryMatch {
			return false
		}
	}

	return true
}

func (h *OrderHandler) extractOrderDetails(c *gin.Context, estimateRequest dto.EstimateRequest) ([]dto.OrderDetail, error) {
	var orderDetails []dto.OrderDetail

	for _, order := range estimateRequest.Orders {
		// Get merchant details from database
		merchantUUID, err := uuid.Parse(order.MerchantId)
		if err != nil {
			continue // Skip invalid merchant ID
		}

		merchantPgUUID := pgtype.UUID{Bytes: merchantUUID, Valid: true}
		merchant, err := h.Q.GetMerchantDetailsByID(c, merchantPgUUID)
		if err != nil {
			continue // Skip if merchant not found
		}

		// Get items details from database
		var orderItems []dto.OrderItem
		for _, item := range order.Items {
			itemUUID, err := uuid.Parse(item.ItemId)
			if err != nil {
				continue // Skip invalid item ID
			}

			itemPgUUID := pgtype.UUID{Bytes: itemUUID, Valid: true}
			merchantItem, err := h.Q.GetMerchantItemByID(c, itemPgUUID)
			if err != nil {
				continue // Skip if item not found
			}

			orderItem := dto.OrderItem{
				ItemID:          item.ItemId,
				Name:            merchantItem.Name,
				ProductCategory: string(merchantItem.ProductCategory),
				Price:           int(merchantItem.Price),
				Quantity:        item.Quantity,
				ImageURL:        merchantItem.ImageUrl,
				CreatedAt:       merchantItem.CreatedAt.Time.Format("2006-01-02T15:04:05.000000000Z07:00"),
			}
			orderItems = append(orderItems, orderItem)
		}

		lat64, _ := merchant.Lat.(float64)
		long64, _ := merchant.Long.(float64)

		orderDetail := dto.OrderDetail{
			Merchant: dto.OrderMerchant{
				MerchantID:       order.MerchantId,
				Name:             merchant.Name,
				MerchantCategory: string(merchant.MerchantCategory),
				ImageURL:         merchant.ImageUrl,
				Location: dto.OrderLocation{
					Lat:  lat64,
					Long: long64,
				},
				CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02T15:04:05.000000000Z07:00"),
			},
			Items: orderItems,
		}
		orderDetails = append(orderDetails, orderDetail)
	}

	return orderDetails, nil
}
