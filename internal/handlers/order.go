package handlers

import (
	"encoding/json"
	"net/http"

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
	response := h.buildOrdersResponse(c, orders, params)

	c.JSON(http.StatusOK, response)
}

func (h *OrderHandler) buildOrdersResponse(c *gin.Context, orders []db.GetOrdersByUserIDRow, params dto.GetOrdersParams) dto.GetOrdersResponse {
	var response dto.GetOrdersResponse
	filteredCount := 0

	for _, order := range orders {
		// Stop if we have enough filtered results
		if filteredCount >= params.Limit {
			break
		}

		orderIDStr := order.ID.String()

		// Parse estimate data
		var estimateRequest dto.EstimateRequest
		if err := json.Unmarshal(order.EstimateData, &estimateRequest); err != nil {
			continue
		}

		// Apply filters
		if !h.matchesFilters(estimateRequest, params) {
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

		response = append(response, orderHistory)
		filteredCount++
	}

	return response
}

func (h *OrderHandler) matchesFilters(estimateRequest dto.EstimateRequest, params dto.GetOrdersParams) bool {
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

	// For name and merchant category filters, we need to query the database
	// This would require additional queries to merchants and items
	// For now, we'll return true for these filters
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
