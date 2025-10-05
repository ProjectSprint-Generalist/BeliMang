package handlers

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sync"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/db"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/shared"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EstimateHandler struct {
	Q *db.Queries
}

func NewEstimateHandler(pool *pgxpool.Pool) *EstimateHandler {
	q := db.New(pool)
	return &EstimateHandler{Q: q}
}

func (h *EstimateHandler) Estimate(c *gin.Context) {
	var req dto.EstimateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
			Code:    http.StatusBadRequest,
		})
		return
	}

	startCount := 0
	for _, o := range req.Orders {
		if o.IsStartingPoint {
			startCount++
		}
	}
	if startCount != 1 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "There must be exactly one starting point",
			Code:    http.StatusBadRequest,
		})
		return
	}

	ctx := c.Request.Context()
	totalPrice := 0.0
	maxDistance := 0.0

	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(req.Orders))
	concurrency := make(chan struct{}, 4)

	for _, order := range req.Orders {
		wg.Add(1)
		go func(order dto.EstimateOrder) {
			defer wg.Done()
			concurrency <- struct{}{}
			defer func() { <-concurrency }()

			merchant, err := h.Q.GetMerchantLocationByID(ctx, order.MerchantId)
			if err != nil {
				errChan <- err
				return
			}

			dist := shared.Haversine(
				req.UserLocation.Lat,
				req.UserLocation.Long,
				merchant.Lat,
				merchant.Long,
			)
			log.Printf("Distance to merchant %s: %.2f km", order.MerchantId, dist)

			if dist > 3 {
				errChan <- &DistanceError{}
				return
			}

			mu.Lock()
			if dist > maxDistance {
				maxDistance = dist
			}
			mu.Unlock()

			for _, item := range order.Items {
				itemData, err := h.Q.GetMerchantItemPriceByID(ctx, item.ItemId)
				if err != nil {
					errChan <- err
					return
				}

				mu.Lock()
				totalPrice += float64(itemData.Price * int32(item.Quantity))
				mu.Unlock()
			}
		}(order)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		switch err.(type) {
		case *DistanceError:
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Error:   "One of the merchants is too far (over 3km)",
				Code:    http.StatusBadRequest,
			})
			return
		default:
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Error:   "One of the merchants or items not found",
				Code:    http.StatusNotFound,
			})
			return
		}
	}

	const speed = 40.0
	deliveryTime := (maxDistance / speed) * 60

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
	user, err := h.Q.GetUserByUsername(ctx, username.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Error:   "User not found",
			Code:    http.StatusUnauthorized,
		})
		return
	}

	// Round and store values
	roundedTotalPrice := math.Round(totalPrice*100) / 100
	roundedDeliveryTime := math.Round(deliveryTime*100) / 100

	// Save estimate data to JSON
	rawJSON, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to store estimate data",
			Code:    http.StatusInternalServerError,
		})
		return
	}
	estimateData := rawJSON

	// Store the calculated estimate in the database
	estimateID, err := h.Q.CreateCalculatedEstimate(ctx, db.CreateCalculatedEstimateParams{
		UserID:                       user.ID,
		TotalPrice:                   int32(roundedTotalPrice),
		EstimatedDeliveryTimeMinutes: int32(roundedDeliveryTime),
		EstimateData:                 estimateData,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "Failed to store estimate",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, dto.EstimateResponse{
		TotalPrice:                  roundedTotalPrice,
		EstimatedDeliveryTimeInMins: roundedDeliveryTime,
		CalculatedEstimateID:        estimateID.String(),
	})
}

type DistanceError struct{}

func (e *DistanceError) Error() string { return "distance exceeds 3km" }
