package repository

import (
	"context"
	"errors"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/db"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/merchant/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MerchantRepositoryImpl struct {
	db *pgxpool.Pool
	q  *db.Queries
}

func NewMerchantRepositoryImpl(pool *pgxpool.Pool) MerchantRepository {
	return &MerchantRepositoryImpl{
		db: pool,
		q:  db.New(pool),
	}
}

func (r *MerchantRepositoryImpl) CreateMerchant(ctx context.Context, merchant domain.Merchant) (string, error) {
	id, err := r.q.CreateMerchant(ctx, db.CreateMerchantParams{
		Name:             merchant.Name,
		MerchantCategory: db.MerchantCategory(merchant.MerchantCategory),
		ImageUrl:         merchant.ImageURL,
		StMakepoint:      merchant.Location.Long,
		StMakepoint_2:    merchant.Location.Lat,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return "", ErrMerchantAlreadyExists
			}
			return "", ErrInternalServerError
		}
		return "", ErrInternalServerError
	}
	return id.String(), nil
}

func (r *MerchantRepositoryImpl) GetMerchants(ctx context.Context, merchantID, name, category, createdAt string, limit, offset int32) ([]domain.Merchant, int64, error) {
	// Validate and set default for createdAt parameter
	if createdAt != "" && createdAt != "asc" && createdAt != "desc" {
		createdAt = "desc" // default to desc if invalid value provided
	}
	if createdAt == "" {
		createdAt = "desc" // default to desc if not provided
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
	if category != "" {
		categoryText = pgtype.Text{String: category, Valid: true}
	}

	// Get total count
	total, err := r.q.CountMerchants(ctx, db.CountMerchantsParams{
		MerchantID:       merchantIDText,
		MerchantCategory: categoryText,
		Name:             nameText,
	})
	if err != nil {
		return nil, 0, ErrInternalServerError
	}

	// Get merchants
	merchants, err := r.q.GetMerchants(ctx, db.GetMerchantsParams{
		MerchantID:       merchantIDText,
		MerchantCategory: categoryText,
		Name:             nameText,
		CreatedAt:        createdAt,
		OffsetVal:        offset,
		LimitVal:         limit,
	})
	if err != nil {
		return nil, 0, ErrInternalServerError
	}

	// Convert to domain objects
	domainMerchants := make([]domain.Merchant, 0, len(merchants))
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

		domainMerchants = append(domainMerchants, domain.Merchant{
			ID:               uuidStr,
			Name:             m.Name,
			MerchantCategory: domain.MerchantCategory(m.MerchantCategory),
			Location: domain.Location{
				Lat:  lat,
				Long: long,
			},
			ImageURL:  m.ImageUrl,
			CreatedAt: m.CreatedAt.Time.Format("2006-01-02T15:04:05.000Z"),
		})
	}

	return domainMerchants, total, nil
}

// GetMerchantByID checks if a merchant exists by ID
func (r *MerchantRepositoryImpl) GetMerchantByID(ctx context.Context, merchantID string) (bool, error) {
	var merchantUUID pgtype.UUID
	if err := merchantUUID.Scan(merchantID); err != nil {
		return false, ErrInvalidUUID
	}

	exists, err := r.q.GetMerchantByID(ctx, merchantUUID)
	if err != nil {
		return false, ErrInternalServerError
	}

	return exists, nil
}

// CountMerchants counts merchants matching the given filters
func (r *MerchantRepositoryImpl) CountMerchants(ctx context.Context, merchantID, name, category string) (int64, error) {
	var merchantIDText pgtype.Text
	if merchantID != "" {
		merchantIDText = pgtype.Text{String: merchantID, Valid: true}
	}
	var nameText pgtype.Text
	if name != "" {
		nameText = pgtype.Text{String: name, Valid: true}
	}
	var categoryText pgtype.Text
	if category != "" {
		categoryText = pgtype.Text{String: category, Valid: true}
	}

	total, err := r.q.CountMerchants(ctx, db.CountMerchantsParams{
		MerchantID:       merchantIDText,
		MerchantCategory: categoryText,
		Name:             nameText,
	})
	if err != nil {
		return 0, ErrInternalServerError
	}

	return total, nil
}

// CreateMerchantItem creates a new merchant item
func (r *MerchantRepositoryImpl) CreateMerchantItem(ctx context.Context, item domain.MerchantItem) (string, error) {
	var merchantUUID pgtype.UUID
	if err := merchantUUID.Scan(item.MerchantID); err != nil {
		return "", ErrInvalidUUID
	}

	id, err := r.q.CreateMerchantItem(ctx, db.CreateMerchantItemParams{
		MerchantID:      merchantUUID,
		Name:            item.Name,
		ProductCategory: db.ProductCategory(item.ProductCategory),
		Price:           item.Price,
		ImageUrl:        item.ImageURL,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return "", ErrMerchantItemNotFound // or appropriate error
			}
			return "", ErrInternalServerError
		}
		return "", ErrInternalServerError
	}
	return id.String(), nil
}

// GetMerchantItems retrieves merchant items with filtering and pagination
func (r *MerchantRepositoryImpl) GetMerchantItems(ctx context.Context, merchantID, itemID, name, category, createdAt string, limit, offset int32) ([]domain.MerchantItem, int64, error) {
	// Validate and set default for createdAt parameter
	if createdAt != "" && createdAt != "asc" && createdAt != "desc" {
		createdAt = "desc" // default to desc if invalid value provided
	}
	if createdAt == "" {
		createdAt = "desc" // default to desc if not provided
	}

	var merchantUUID pgtype.UUID
	if err := merchantUUID.Scan(merchantID); err != nil {
		return nil, 0, ErrInvalidUUID
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
	if category != "" {
		productCategoryText = pgtype.Text{String: category, Valid: true}
	}

	// Get total count
	total, err := r.q.CountMerchantItems(ctx, db.CountMerchantItemsParams{
		MerchantID:      merchantUUID,
		ItemID:          itemIDText,
		ProductCategory: productCategoryText,
		Name:            nameText,
	})
	if err != nil {
		return nil, 0, ErrInternalServerError
	}

	// Get merchant items
	items, err := r.q.GetMerchantItems(ctx, db.GetMerchantItemsParams{
		MerchantID:      merchantUUID,
		ItemID:          itemIDText,
		ProductCategory: productCategoryText,
		Name:            nameText,
		CreatedAt:       createdAt,
		OffsetVal:       offset,
		LimitVal:        limit,
	})
	if err != nil {
		return nil, 0, ErrInternalServerError
	}

	// Convert to domain objects
	domainItems := make([]domain.MerchantItem, 0, len(items))
	for _, item := range items {
		// Convert UUID bytes to string
		uuidStr := ""
		if item.ID.Valid {
			uuidStr = pgtype.UUID{Bytes: item.ID.Bytes, Valid: true}.String()
		}

		domainItems = append(domainItems, domain.MerchantItem{
			ID:              uuidStr,
			MerchantID:      merchantID,
			Name:            item.Name,
			ProductCategory: domain.ProductCategory(item.ProductCategory),
			Price:           item.Price,
			ImageURL:        item.ImageUrl,
			CreatedAt:       item.CreatedAt.Time.Format("2006-01-02T15:04:05.000Z"),
		})
	}

	return domainItems, total, nil
}

// CountMerchantItems counts merchant items matching the given filters
func (r *MerchantRepositoryImpl) CountMerchantItems(ctx context.Context, merchantID, itemID, name, category string) (int64, error) {
	var merchantUUID pgtype.UUID
	if err := merchantUUID.Scan(merchantID); err != nil {
		return 0, ErrInvalidUUID
	}

	var itemIDText pgtype.Text
	if itemID != "" {
		itemIDText = pgtype.Text{String: itemID, Valid: true}
	}
	var nameText pgtype.Text
	if name != "" {
		nameText = pgtype.Text{String: name, Valid: true}
	}
	var productCategoryText pgtype.Text
	if category != "" {
		productCategoryText = pgtype.Text{String: category, Valid: true}
	}

	total, err := r.q.CountMerchantItems(ctx, db.CountMerchantItemsParams{
		MerchantID:      merchantUUID,
		ItemID:          itemIDText,
		ProductCategory: productCategoryText,
		Name:            nameText,
	})
	if err != nil {
		return 0, ErrInternalServerError
	}

	return total, nil
}

// GetNearbyMerchants retrieves nearby merchants based on coordinates
func (r *MerchantRepositoryImpl) GetNearbyMerchants(ctx context.Context, lat, long float64, merchantID, name, category string, limit, offset int32) ([]domain.Merchant, int64, error) {
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
	if category != "" {
		categoryText = pgtype.Text{String: category, Valid: true}
	}

	// Get total count
	total, err := r.q.CountNearbyMerchants(ctx, db.CountNearbyMerchantsParams{
		MerchantID:       merchantIDText,
		MerchantCategory: categoryText,
		Name:             nameText,
	})
	if err != nil {
		return nil, 0, ErrInternalServerError
	}

	// Get nearby merchants
	rows, err := r.q.GetNearbyMerchants(ctx, db.GetNearbyMerchantsParams{
		Lat:              lat,
		Long:             long,
		MerchantID:       merchantIDText,
		MerchantCategory: categoryText,
		Name:             nameText,
		RowLimit:         limit,
		RowOffset:        offset,
	})
	if err != nil {
		return nil, 0, ErrInternalServerError
	}

	// Convert to domain objects
	domainMerchants := make([]domain.Merchant, 0, len(rows))
	for _, m := range rows {
		// Convert UUID to string
		merchantIDStr := ""
		if m.ID.Valid {
			merchantIDStr = pgtype.UUID{Bytes: m.ID.Bytes, Valid: true}.String()
		}

		lat64, _ := m.Lat.(float64)
		long64, _ := m.Long.(float64)

		domainMerchants = append(domainMerchants, domain.Merchant{
			ID:               merchantIDStr,
			Name:             m.Name,
			MerchantCategory: domain.MerchantCategory(m.MerchantCategory),
			Location: domain.Location{
				Lat:  lat64,
				Long: long64,
			},
			ImageURL:  m.ImageUrl,
			CreatedAt: m.CreatedAt.Time.Format("2006-01-02T15:04:05.000Z"),
		})
	}

	return domainMerchants, total, nil
}

// CountNearbyMerchants counts nearby merchants matching the given filters
func (r *MerchantRepositoryImpl) CountNearbyMerchants(ctx context.Context, lat, long float64, merchantID, name, category string) (int64, error) {
	var merchantIDText pgtype.Text
	if merchantID != "" {
		merchantIDText = pgtype.Text{String: merchantID, Valid: true}
	}
	var nameText pgtype.Text
	if name != "" {
		nameText = pgtype.Text{String: name, Valid: true}
	}
	var categoryText pgtype.Text
	if category != "" {
		categoryText = pgtype.Text{String: category, Valid: true}
	}

	total, err := r.q.CountNearbyMerchants(ctx, db.CountNearbyMerchantsParams{
		MerchantID:       merchantIDText,
		MerchantCategory: categoryText,
		Name:             nameText,
	})
	if err != nil {
		return 0, ErrInternalServerError
	}

	return total, nil
}
