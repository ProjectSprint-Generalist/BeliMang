package handler

import (
	"net/http"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/order/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// OrderHandler handles order-related HTTP requests
type OrderHandler struct {
	svc *service.OrderService
}

// NewOrderHandler creates a new order handler
func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// Estimate handles POST /users/estimate
func (h *OrderHandler) Estimate(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var req dto.EstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
			Code:    http.StatusBadRequest,
		})
		return
	}

	resp, err := h.svc.CalculateEstimate(c.Request.Context(), userID.(string), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case service.ErrInvalidEstimateRequest, service.ErrInvalidStartingPoint, service.ErrCoordinatesTooFar:
			statusCode = http.StatusBadRequest
		case service.ErrMerchantNotFound, service.ErrMerchantItemNotFound:
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    statusCode,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// CreateOrder handles POST /users/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
			Code:    http.StatusBadRequest,
		})
		return
	}

	resp, err := h.svc.CreateOrder(c.Request.Context(), userID.(string), req.CalculatedEstimateID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case service.ErrCalculatedEstimateNotFound, service.ErrUnauthorizedAccess, service.ErrInvalidOrderRequest:
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    statusCode,
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetOrders handles GET /users/orders
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Error:   "Unauthorized",
			Code:    http.StatusUnauthorized,
		})
		return
	}

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

	resp, err := h.svc.GetOrderHistory(c.Request.Context(), userID.(pgtype.UUID), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
