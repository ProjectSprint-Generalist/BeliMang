package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/db"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/order/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrderRepositoryImpl implements OrderRepository using PostgreSQL
type OrderRepositoryImpl struct {
	db *pgxpool.Pool
	q  *db.Queries
}

// NewOrderRepositoryImpl creates a new order repository implementation
func NewOrderRepositoryImpl(pool *pgxpool.Pool) OrderRepository {
	return &OrderRepositoryImpl{
		db: pool,
		q:  db.New(pool),
	}
}

// CreateCalculatedEstimate creates a new calculated estimate in the database
func (r *OrderRepositoryImpl) CreateCalculatedEstimate(ctx context.Context, userID pgtype.UUID, totalPrice int32, deliveryTimeMinutes int32, estimateData []byte) (string, error) {
	id, err := r.q.CreateCalculatedEstimate(ctx, db.CreateCalculatedEstimateParams{
		UserID:                       userID,
		TotalPrice:                   totalPrice,
		EstimatedDeliveryTimeMinutes: deliveryTimeMinutes,
		EstimateData:                 estimateData,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return "", ErrCalculatedEstimateAlreadyExists
			}
			return "", ErrInternalServerError
		}
		return "", ErrInternalServerError
	}
	return id.String(), nil
}

// GetCalculatedEstimateByID retrieves a calculated estimate by ID
func (r *OrderRepositoryImpl) GetCalculatedEstimateByID(ctx context.Context, estimateID pgtype.UUID) (domain.CalculatedEstimate, error) {
	var estimateUUID pgtype.UUID
	if err := estimateUUID.Scan(estimateID); err != nil {
		return domain.CalculatedEstimate{}, ErrInvalidUUID
	}

	estimate, err := r.q.GetCalculatedEstimateByID(ctx, estimateUUID)
	if err != nil {
		return domain.CalculatedEstimate{}, ErrCalculatedEstimateNotFound
	}

	// Parse estimate data JSON
	var estimateRequest dto.EstimateRequest
	if err := json.Unmarshal(estimate.EstimateData, &estimateRequest); err != nil {
		return domain.CalculatedEstimate{}, ErrInternalServerError
	}

	return domain.CalculatedEstimate{
		ID:                           estimate.ID.String(),
		UserID:                       estimate.UserID.String(),
		TotalPrice:                   estimate.TotalPrice,
		EstimatedDeliveryTimeMinutes: estimate.EstimatedDeliveryTimeMinutes,
		EstimateData:                 estimate.EstimateData,
		CreatedAt:                    estimate.CreatedAt.Time,
	}, nil
}

// CreateOrder creates a new order in the database
func (r *OrderRepositoryImpl) CreateOrder(ctx context.Context, order domain.Order) (pgtype.UUID, error) {
	var userUUID pgtype.UUID
	if err := userUUID.Scan(order.UserID); err != nil {
		return pgtype.UUID{}, ErrInvalidUUID
	}

	var estimateUUID pgtype.UUID
	if err := estimateUUID.Scan(order.CalculatedEstimateID); err != nil {
		return pgtype.UUID{}, ErrInvalidUUID
	}

	id, err := r.q.CreateOrder(ctx, db.CreateOrderParams{
		UserID:               userUUID,
		CalculatedEstimateID: estimateUUID,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return pgtype.UUID{}, ErrOrderAlreadyExists
			}
			return pgtype.UUID{}, ErrInternalServerError
		}
		return pgtype.UUID{}, ErrInternalServerError
	}
	return id, nil
}

// GetOrdersByUserID retrieves orders for a user with pagination
func (r *OrderRepositoryImpl) GetOrdersByUserID(ctx context.Context, userID pgtype.UUID, limit, offset int32) ([]domain.OrderWithEstimate, error) {
	orders, err := r.q.GetOrdersByUserID(ctx, db.GetOrdersByUserIDParams{
		Column1: userID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, ErrInternalServerError
	}

	// Convert to domain objects
	domainOrders := make([]domain.OrderWithEstimate, 0, len(orders))
	for _, order := range orders {
		// Parse estimate data JSON
		var estimateRequest dto.EstimateRequest
		if err := json.Unmarshal(order.EstimateData, &estimateRequest); err != nil {
			continue // Skip orders with invalid estimate data
		}

		domainOrders = append(domainOrders, domain.OrderWithEstimate{
			ID:           order.ID.String(),
			UserID:       order.UserID.String(),
			EstimateData: estimateRequest,
			CreatedAt:    order.CreatedAt.Time,
		})
	}

	return domainOrders, nil
}

// GetOrdersCountByUserID gets the total count of orders for a user
func (r *OrderRepositoryImpl) GetOrdersCountByUserID(ctx context.Context, userID pgtype.UUID) (int64, error) {
	count, err := r.q.GetOrdersCountByUserID(ctx, userID)
	if err != nil {
		return 0, ErrInternalServerError
	}

	return count, nil
}

// GetFilteredOrdersCount gets the count of orders for a user with filtering applied
func (r *OrderRepositoryImpl) GetFilteredOrdersCount(ctx context.Context, userID pgtype.UUID, merchantID, name, category *string) (int64, error) {
	// For now, just return the total count since filtering is done client-side
	// In the future, this could be implemented with server-side filtering
	return r.GetOrdersCountByUserID(ctx, userID)
}

// GetMerchantDetailsByID retrieves merchant details for order history
func (r *OrderRepositoryImpl) GetMerchantDetailsByID(ctx context.Context, merchantID pgtype.UUID) (dto.OrderMerchant, error) {
	merchant, err := r.q.GetMerchantDetailsByID(ctx, merchantID)
	if err != nil {
		return dto.OrderMerchant{}, ErrMerchantNotFound
	}

	// Handle lat/long conversion from interface{} to float64
	var lat, long float64
	if merchant.Lat != nil {
		if latVal, ok := merchant.Lat.(float64); ok {
			lat = latVal
		}
	}
	if merchant.Long != nil {
		if longVal, ok := merchant.Long.(float64); ok {
			long = longVal
		}
	}

	return dto.OrderMerchant{
		MerchantID:       merchantID,
		Name:             merchant.Name,
		MerchantCategory: string(merchant.MerchantCategory),
		ImageURL:         merchant.ImageUrl,
		Location: dto.OrderLocation{
			Lat:  lat,
			Long: long,
		},
		CreatedAt: merchant.CreatedAt.Time.Format("2006-01-02T15:04:05.000000000Z07:00"),
	}, nil
}

// GetMerchantItemByID retrieves merchant item details for order history
func (r *OrderRepositoryImpl) GetMerchantItemByID(ctx context.Context, itemID pgtype.UUID) (dto.OrderItem, error) {
	item, err := r.q.GetMerchantItemByID(ctx, itemID)
	if err != nil {
		return dto.OrderItem{}, ErrMerchantItemNotFound
	}

	return dto.OrderItem{
		ItemID:          itemID,
		Name:            item.Name,
		ProductCategory: string(item.ProductCategory),
		Price:           int(item.Price),
		ImageURL:        item.ImageUrl,
		CreatedAt:       item.CreatedAt.Time.Format("2006-01-02T15:04:05.000000000Z07:00"),
	}, nil
}

// GetMerchantLocationByID retrieves merchant location for distance calculations
func (r *OrderRepositoryImpl) GetMerchantLocationByID(ctx context.Context, merchantID pgtype.UUID) (dto.Location, error) {
	location, err := r.q.GetMerchantLocationByID(ctx, merchantID)
	if err != nil {
		return dto.Location{}, ErrMerchantNotFound
	}

	return dto.Location{
		Lat:  location.Lat,
		Long: location.Long,
	}, nil
}

// GetMerchantItemPriceByID retrieves item price for price calculations
func (r *OrderRepositoryImpl) GetMerchantItemPriceByID(ctx context.Context, itemID pgtype.UUID) (int32, error) {
	item, err := r.q.GetMerchantItemPriceByID(ctx, itemID)
	if err != nil {
		return 0, ErrMerchantItemNotFound
	}

	return item.Price, nil
}

// GetMerchantItemsByMerchantID retrieves all items for a merchant (optimized for order history)
func (r *OrderRepositoryImpl) GetMerchantItemsByMerchantID(ctx context.Context, merchantID pgtype.UUID) ([]dto.OrderItem, error) {
	items, err := r.q.GetMerchantItems(ctx, db.GetMerchantItemsParams{
		MerchantID: merchantID,
		// Empty filters to get all items for this merchant
		ItemID:          pgtype.Text{},
		ProductCategory: pgtype.Text{},
		Name:            pgtype.Text{},
		CreatedAt:       "desc",
		OffsetVal:       0,
		LimitVal:        1000, // Get all items for this merchant
	})
	if err != nil {
		return nil, ErrInternalServerError
	}

	// Convert to DTOs
	dtoItems := make([]dto.OrderItem, 0, len(items))
	for _, item := range items {
		// Convert UUID bytes to string
		itemID := ""
		if item.ID.Valid {
			itemID = pgtype.UUID{Bytes: item.ID.Bytes, Valid: true}.String()
		}

		dtoItems = append(dtoItems, dto.OrderItem{
			ItemID:          itemID,
			Name:            item.Name,
			ProductCategory: string(item.ProductCategory),
			Price:           int(item.Price),
			ImageURL:        item.ImageUrl,
			CreatedAt:       item.CreatedAt.Time.Format("2006-01-02T15:04:05.000000000Z07:00"),
		})
	}

	return dtoItems, nil
}
