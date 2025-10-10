package service

import (
	"context"
	"errors"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/merchant/domain"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/merchant/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/merchant/repository"
)

var (
	ErrInvalidMerchantData     = errors.New("invalid merchant data")
	ErrInvalidMerchantItemData = errors.New("invalid merchant item data")
	ErrMerchantNotFound        = errors.New("merchant not found")
	ErrInvalidPagination       = errors.New("invalid pagination parameters")
)

type MerchantService struct {
	repo repository.MerchantRepository
}

func NewMerchantService(repo repository.MerchantRepository) *MerchantService {
	return &MerchantService{repo: repo}
}

func (s *MerchantService) CreateMerchant(ctx context.Context, req dto.MerchantCreateRequest) (dto.MerchantCreateResponse, error) {
	merchant, err := domain.NewMerchant(
		req.Name,
		domain.MerchantCategory(req.MerchantCategory),
		domain.Location{Lat: req.Location.Lat, Long: req.Location.Long},
		req.ImageURL,
	)
	if err != nil {
		return dto.MerchantCreateResponse{}, ErrInvalidMerchantData
	}

	id, err := s.repo.CreateMerchant(ctx, merchant)
	if err != nil {
		if errors.Is(err, repository.ErrMerchantAlreadyExists) {
			return dto.MerchantCreateResponse{}, ErrInvalidMerchantData
		}
		return dto.MerchantCreateResponse{}, err
	}

	return dto.MerchantCreateResponse{MerchantId: id}, nil
}

func (s *MerchantService) GetMerchants(ctx context.Context, merchantID, name, category, createdAt string, limit, offset int32) (dto.GetMerchantsResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 5 // default
	}
	if offset < 0 {
		offset = 0
	}

	merchants, total, err := s.repo.GetMerchants(ctx, merchantID, name, category, createdAt, limit, offset)
	if err != nil {
		return dto.GetMerchantsResponse{}, err
	}

	merchantData := make([]dto.MerchantData, 0, len(merchants))
	for _, m := range merchants {
		merchantData = append(merchantData, dto.ConvertDomainToDTO(m))
	}

	return dto.GetMerchantsResponse{
		Data: merchantData,
		Meta: dto.MerchantMeta{
			Limit:  int(limit),
			Offset: int(offset),
			Total:  int(total),
		},
	}, nil
}

func (s *MerchantService) CreateMerchantItem(ctx context.Context, merchantID string, req dto.MerchantItemCreateRequest) (dto.MerchantItemCreateResponse, error) {
	exists, err := s.repo.GetMerchantByID(ctx, merchantID)
	if err != nil {
		return dto.MerchantItemCreateResponse{}, err
	}
	if !exists {
		return dto.MerchantItemCreateResponse{}, ErrMerchantNotFound
	}

	item, err := domain.NewMerchantItem(
		req.Name,
		domain.ProductCategory(req.ProductCategory),
		int32(req.Price),
		req.ImageURL,
	)
	if err != nil {
		return dto.MerchantItemCreateResponse{}, ErrInvalidMerchantItemData
	}

	item.MerchantID = merchantID

	id, err := s.repo.CreateMerchantItem(ctx, item)
	if err != nil {
		return dto.MerchantItemCreateResponse{}, err
	}

	return dto.MerchantItemCreateResponse{ItemId: id}, nil
}

func (s *MerchantService) GetMerchantItems(ctx context.Context, merchantID, itemID, name, category, createdAt string, limit, offset int32) (dto.GetMerchantItemsResponse, error) {
	exists, err := s.repo.GetMerchantByID(ctx, merchantID)
	if err != nil {
		return dto.GetMerchantItemsResponse{}, err
	}
	if !exists {
		return dto.GetMerchantItemsResponse{}, ErrMerchantNotFound
	}

	if limit <= 0 || limit > 100 {
		limit = 5 // default
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := s.repo.GetMerchantItems(ctx, merchantID, itemID, name, category, createdAt, limit, offset)
	if err != nil {
		return dto.GetMerchantItemsResponse{}, err
	}

	itemData := make([]dto.MerchantItemData, 0, len(items))
	for _, item := range items {
		itemData = append(itemData, dto.ConvertDomainItemToDTO(item))
	}

	return dto.GetMerchantItemsResponse{
		Data: itemData,
		Meta: dto.MerchantMeta{
			Limit:  int(limit),
			Offset: int(offset),
			Total:  int(total),
		},
	}, nil
}

func (s *MerchantService) GetNearbyMerchants(ctx context.Context, lat, long float64, merchantID, name, category string, limit, offset int32) (dto.GetNearbyMerchantsResponse, error) {
	if lat < -90 || lat > 90 || long < -180 || long > 180 {
		return dto.GetNearbyMerchantsResponse{}, ErrInvalidMerchantData
	}

	if limit <= 0 || limit > 100 {
		limit = 5 // default
	}
	if offset < 0 {
		offset = 0
	}

	merchants, total, err := s.repo.GetNearbyMerchants(ctx, lat, long, merchantID, name, category, limit, offset)
	if err != nil {
		return dto.GetNearbyMerchantsResponse{}, err
	}

	nearbyMerchants := make([]dto.NearbyMerchant, 0, len(merchants))
	for _, m := range merchants {
		items, _, err := s.repo.GetMerchantItems(ctx, m.ID, "", "", "", "desc", limit, offset) // Get all items for this merchant
		if err != nil {
			continue // Skip this merchant if we can't get items
		}

		itemData := make([]dto.MerchantItemData, 0, len(items))
		for _, item := range items {
			itemData = append(itemData, dto.ConvertDomainItemToDTO(item))
		}

		nearbyMerchants = append(nearbyMerchants, dto.NearbyMerchant{
			Merchant: dto.ConvertDomainToDTO(m),
			Items:    itemData,
		})
	}

	return dto.GetNearbyMerchantsResponse{
		Data: nearbyMerchants,
		Meta: dto.MerchantMeta{
			Limit:  int(limit),
			Offset: int(offset),
			Total:  int(total),
		},
	}, nil
}
