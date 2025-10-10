package repository

import (
	"context"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/merchant/domain"
)

type MerchantRepository interface {
	// Merchant operations
	CreateMerchant(ctx context.Context, merchant domain.Merchant) (string, error)
	GetMerchants(ctx context.Context, merchantID, name, category, createdAt string, limit, offset int32) ([]domain.Merchant, int64, error)
	GetMerchantByID(ctx context.Context, merchantID string) (bool, error)
	CountMerchants(ctx context.Context, merchantID, name, category string) (int64, error)

	// Merchant item operations
	CreateMerchantItem(ctx context.Context, item domain.MerchantItem) (string, error)
	GetMerchantItems(ctx context.Context, merchantID, itemID, name, category, createdAt string, limit, offset int32) ([]domain.MerchantItem, int64, error)
	CountMerchantItems(ctx context.Context, merchantID, itemID, name, category string) (int64, error)

	// Nearby merchant operations
	GetNearbyMerchants(ctx context.Context, lat, long float64, merchantID, name, category string, limit, offset int32) ([]domain.Merchant, int64, error)
	CountNearbyMerchants(ctx context.Context, lat, long float64, merchantID, name, category string) (int64, error)
}
