package service

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strings"
	"sync"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/order/domain"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/order/repository"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/shared"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	DELIVERY_SPEED_KMH = 40.0
)

var (
	ErrInvalidEstimateRequest     = errors.New("invalid estimate request")
	ErrInvalidStartingPoint       = errors.New("invalid starting point configuration")
	ErrCoordinatesTooFar          = errors.New("coordinates are too far from merchants")
	ErrMerchantNotFound           = errors.New("merchant not found")
	ErrMerchantItemNotFound       = errors.New("merchant item not found")
	ErrInvalidOrderRequest        = errors.New("invalid order request")
	ErrCalculatedEstimateNotFound = errors.New("calculated estimate not found")
	ErrUnauthorizedAccess         = errors.New("unauthorized access to estimate")
	ErrInvalidPagination          = errors.New("invalid pagination parameters")
)

// OrderService handles order business logic
type OrderService struct {
	repo repository.OrderRepository
}

// NewOrderService creates a new order service
func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// CalculateEstimate calculates price and delivery time for an order
func (s *OrderService) CalculateEstimate(ctx context.Context, userID string, req dto.EstimateRequest) (dto.EstimateResponse, error) {
	// Validate request structure
	if err := s.validateEstimateRequest(req); err != nil {
		return dto.EstimateResponse{}, err
	}

	// Validate coordinates distance and get max distance (concurrent)
	maxDistance, err := s.validateCoordinatesDistance(ctx, dto.Location{Lat: req.UserLocation.Lat, Long: req.UserLocation.Long}, req.Orders)
	if err != nil {
		return dto.EstimateResponse{}, err
	}

	// Calculate total price (concurrent)
	totalPrice, err := s.calculateTotalPrice(ctx, req.Orders)
	if err != nil {
		return dto.EstimateResponse{}, err
	}

	deliveryTime := (maxDistance / DELIVERY_SPEED_KMH) * 60
	estimatedTime := int32(math.Round(deliveryTime))

	roundedTotalPrice := int32(math.Round(float64(totalPrice)*100) / 100)

	// Save estimate data
	estimateData, err := json.Marshal(req)
	if err != nil {
		return dto.EstimateResponse{}, ErrInvalidEstimateRequest
	}

	calculatedEstimateID, err := s.repo.CreateCalculatedEstimate(ctx, userID, roundedTotalPrice, estimatedTime, estimateData)
	if err != nil {
		return dto.EstimateResponse{}, err
	}

	return dto.EstimateResponse{
		TotalPrice:                  float64(roundedTotalPrice),
		EstimatedDeliveryTimeInMins: float64(estimatedTime),
		CalculatedEstimateID:        calculatedEstimateID,
	}, nil
}

// CreateOrder creates an order from a calculated estimate
func (s *OrderService) CreateOrder(ctx context.Context, userID, calculatedEstimateID string) (dto.CreateOrderResponse, error) {
	// Validate calculated estimate exists and belongs to user
	estimate, err := s.repo.GetCalculatedEstimateByID(ctx, calculatedEstimateID)
	if err != nil {
		if errors.Is(err, repository.ErrCalculatedEstimateNotFound) {
			return dto.CreateOrderResponse{}, ErrCalculatedEstimateNotFound
		}
		return dto.CreateOrderResponse{}, err
	}

	// Verify ownership
	if estimate.UserID != userID {
		return dto.CreateOrderResponse{}, ErrUnauthorizedAccess
	}

	// Create order domain entity
	order, err := domain.NewOrder(userID, calculatedEstimateID)
	if err != nil {
		return dto.CreateOrderResponse{}, ErrInvalidOrderRequest
	}

	// Save order
	orderID, err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return dto.CreateOrderResponse{}, err
	}

	return dto.CreateOrderResponse{OrderID: orderID}, nil
}

// GetOrderHistory retrieves order history with filtering and pagination
func (s *OrderService) GetOrderHistory(ctx context.Context, userID pgtype.UUID, params dto.GetOrdersParams) (dto.GetOrdersResponse, error) {
	// Get total count first (for all orders, since filtering is done client-side)
	totalCount, err := s.repo.GetOrdersCountByUserID(ctx, userID)
	if err != nil {
		return dto.GetOrdersResponse{}, err
	}

	// Get orders from repository (without filtering - filtering is done client-side)
	orders, err := s.repo.GetOrdersByUserID(ctx, userID, int32(params.Limit*2), 0)
	if err != nil {
		return dto.GetOrdersResponse{}, err
	}

	// Build order history response with filtering
	orderHistoryData, err := s.buildOrderHistoryResponse(ctx, orders, params)
	if err != nil {
		return dto.GetOrdersResponse{}, err
	}

	return dto.GetOrdersResponse{
		Data: orderHistoryData,
		Meta: dto.GetOrdersResponseMeta{
			Limit:  params.Limit,
			Offset: params.Offset,
			Total:  int(totalCount),
		},
	}, nil
}

// validateEstimateRequest validates the estimate request structure
func (s *OrderService) validateEstimateRequest(req dto.EstimateRequest) error {
	if req.UserLocation.Lat < -90 || req.UserLocation.Lat > 90 {
		return ErrInvalidEstimateRequest
	}
	if req.UserLocation.Long < -180 || req.UserLocation.Long > 180 {
		return ErrInvalidEstimateRequest
	}

	if len(req.Orders) == 0 {
		return ErrInvalidEstimateRequest
	}

	// Check starting point validation
	startingPointCount := 0
	for _, order := range req.Orders {
		if order.IsStartingPoint {
			startingPointCount++
		}
	}

	if startingPointCount != 1 {
		return ErrInvalidStartingPoint
	}

	return nil
}

// validateCoordinatesDistance validates that coordinates are within 3km distance and tracks max distance
func (s *OrderService) validateCoordinatesDistance(ctx context.Context, userLocation dto.Location, orders []dto.EstimateOrder) (float64, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(orders))
	concurrency := make(chan struct{}, 4)

	maxDistance := 0.0

	for _, order := range orders {
		wg.Add(1)
		go func(order dto.EstimateOrder) {
			defer wg.Done()
			concurrency <- struct{}{}
			defer func() { <-concurrency }()

			merchant, err := s.repo.GetMerchantLocationByID(ctx, order.MerchantId)
			if err != nil {
				if err == repository.ErrMerchantNotFound {
					errChan <- ErrMerchantNotFound
				} else {
					errChan <- err
				}
				return
			}

			dist := shared.Haversine(
				userLocation.Lat,
				userLocation.Long,
				merchant.Lat,
				merchant.Long,
			)

			mu.Lock()
			if dist > maxDistance {
				maxDistance = dist
			}
			mu.Unlock()

			if dist > 3.0 { // 3km limit from old code
				errChan <- ErrCoordinatesTooFar
				return
			}
		}(order)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != ErrCoordinatesTooFar {
			return 0, err
		}
	}

	return maxDistance, nil
}

// calculateTotalPrice calculates the total price of all items
func (s *OrderService) calculateTotalPrice(ctx context.Context, orders []dto.EstimateOrder) (int32, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(orders))
	concurrency := make(chan struct{}, 4)

	totalPrice := int32(0)

	for _, order := range orders {
		wg.Add(1)
		go func(order dto.EstimateOrder) {
			defer wg.Done()
			concurrency <- struct{}{}
			defer func() { <-concurrency }()

			for _, item := range order.Items {
				itemData, err := s.repo.GetMerchantItemPriceByID(ctx, item.ItemId)
				if err != nil {
					if err == repository.ErrMerchantItemNotFound {
						errChan <- ErrMerchantItemNotFound
					} else {
						errChan <- err
					}
					return
				}

				mu.Lock()
				totalPrice += itemData * int32(item.Quantity)
				mu.Unlock()
			}
		}(order)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		return 0, err
	}

	return totalPrice, nil
}

// buildOrderHistoryResponse builds the order history response with filtering
func (s *OrderService) buildOrderHistoryResponse(ctx context.Context, orders []domain.OrderWithEstimate, params dto.GetOrdersParams) ([]dto.OrderHistory, error) {
	var orderHistory []dto.OrderHistory

	for _, order := range orders {
		// Parse estimate data to get order details
		var orderDetails []dto.OrderDetail

		for _, orderReq := range order.EstimateData.Orders {
			// Get merchant details
			merchant, err := s.repo.GetMerchantDetailsByID(ctx, orderReq.MerchantId)
			if err != nil {
				continue // Skip if merchant not found
			}

			// Apply merchant category filter
			if params.MerchantCategory != nil && string(merchant.MerchantCategory) != *params.MerchantCategory {
				continue
			}

			// Get all items for this merchant (optimized - one query instead of N+1)
			allMerchantItems, err := s.repo.GetMerchantItemsByMerchantID(ctx, orderReq.MerchantId)
			if err != nil {
				continue // Skip if merchant items can't be fetched
			}

			// Filter items based on what's actually ordered and apply name filter
			var orderItems []dto.OrderItem
			for _, item := range allMerchantItems {
				// Check if this item is actually in the order
				found := false
				var quantity int
				for _, itemReq := range orderReq.Items {
					if item.ItemID == itemReq.ItemId {
						found = true
						quantity = itemReq.Quantity
						break
					}
				}

				if !found {
					continue // Skip items not in this order
				}

				// Apply name filter
				if params.Name != nil {
					nameLower := strings.ToLower(*params.Name)
					if !strings.Contains(strings.ToLower(merchant.Name), nameLower) &&
						!strings.Contains(strings.ToLower(item.Name), nameLower) {
						continue
					}
				}

				orderItems = append(orderItems, dto.OrderItem{
					ItemID:          item.ItemID,
					Name:            item.Name,
					ProductCategory: item.ProductCategory,
					Price:           item.Price,
					Quantity:        quantity, // Use the actual quantity from the order
					ImageURL:        item.ImageURL,
					CreatedAt:       item.CreatedAt,
				})
			}

			// Skip if no items match filters
			if len(orderItems) == 0 {
				continue
			}

			orderDetails = append(orderDetails, dto.OrderDetail{
				Merchant: merchant,
				Items:    orderItems,
			})
		}

		// Skip if no order details match filters
		if len(orderDetails) == 0 {
			continue
		}

		orderHistory = append(orderHistory, dto.OrderHistory{
			OrderID: order.ID,
			Orders:  orderDetails,
		})
	}

	return orderHistory, nil
}
