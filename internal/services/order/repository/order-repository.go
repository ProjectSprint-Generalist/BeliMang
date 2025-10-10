package repository

import (
	"context"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/order/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

// OrderRepository defines the interface for order data operations
type OrderRepository interface {
	// Estimate operations
	CreateCalculatedEstimate(ctx context.Context, userID pgtype.UUID, totalPrice int32, deliveryTimeMinutes int32, estimateData []byte) (string, error)
	GetCalculatedEstimateByID(ctx context.Context, estimateID pgtype.UUID) (domain.CalculatedEstimate, error)

	// Order operations
	CreateOrder(ctx context.Context, order domain.Order) (pgtype.UUID, error)
	GetOrdersByUserID(ctx context.Context, userID pgtype.UUID, limit, offset int32) ([]domain.OrderWithEstimate, error)
	GetOrdersCountByUserID(ctx context.Context, userID pgtype.UUID) (int64, error)
	GetFilteredOrdersCount(ctx context.Context, userID pgtype.UUID, merchantID, name, category *string) (int64, error)

	// Merchant and item operations for order history
	GetMerchantDetailsByID(ctx context.Context, merchantID pgtype.UUID) (dto.OrderMerchant, error)
	GetMerchantItemByID(ctx context.Context, itemID pgtype.UUID) (dto.OrderItem, error)

	// Additional operations for estimate calculation
	GetMerchantLocationByID(ctx context.Context, merchantID pgtype.UUID) (dto.Location, error)
	GetMerchantItemPriceByID(ctx context.Context, itemID pgtype.UUID) (int32, error)
	GetMerchantItemsByMerchantID(ctx context.Context, merchantID pgtype.UUID) ([]dto.OrderItem, error)
}
