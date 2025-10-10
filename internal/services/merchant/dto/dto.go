package dto

import "github.com/ProjectSprint-Generalist/BeliMang/internal/services/merchant/domain"

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type MerchantCategory string

const (
	MerchantCategorySmallRestaurant       MerchantCategory = "SmallRestaurant"
	MerchantCategoryMediumRestaurant      MerchantCategory = "MediumRestaurant"
	MerchantCategoryLargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchantCategoryMerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	MerchantCategoryBoothKiosk            MerchantCategory = "BoothKiosk"
	MerchantCategoryConvenienceStore      MerchantCategory = "ConvenienceStore"
)

var ValidMerchantCategories = map[MerchantCategory]bool{
	MerchantCategorySmallRestaurant:       true,
	MerchantCategoryMediumRestaurant:      true,
	MerchantCategoryLargeRestaurant:       true,
	MerchantCategoryMerchandiseRestaurant: true,
	MerchantCategoryBoothKiosk:            true,
	MerchantCategoryConvenienceStore:      true,
}

type ProductCategory string

const (
	ProductCategoryBeverage   ProductCategory = "Beverage"
	ProductCategoryFood       ProductCategory = "Food"
	ProductCategorySnack      ProductCategory = "Snack"
	ProductCategoryCondiments ProductCategory = "Condiments"
	ProductCategoryAdditions  ProductCategory = "Additions"
)

type MerchantCreateRequest struct {
	Name             string           `json:"name" binding:"required,min=2,max=30"`
	MerchantCategory MerchantCategory `json:"merchantCategory" binding:"required"`
	ImageURL         string           `json:"imageURL" binding:"required"`
	Location         Location         `json:"location" binding:"required"`
}

type MerchantCreateResponse struct {
	MerchantId string `json:"merchantId"`
}

type MerchantData struct {
	MerchantID       string   `json:"merchantId"`
	Name             string   `json:"name"`
	MerchantCategory string   `json:"merchantCategory"`
	ImageURL         string   `json:"imageUrl"`
	Location         Location `json:"location"`
	CreatedAt        string   `json:"createdAt"`
}

type MerchantMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type GetMerchantsResponse struct {
	Data []MerchantData `json:"data"`
	Meta MerchantMeta   `json:"meta"`
}

type MerchantItemCreateRequest struct {
	Name            string          `json:"name" binding:"required,min=2,max=30"`
	ProductCategory ProductCategory `json:"productCategory" binding:"required"`
	Price           int             `json:"price" binding:"required,min=1"`
	ImageURL        string          `json:"imageUrl" binding:"required"`
}

type MerchantItemCreateResponse struct {
	ItemId string `json:"itemId"`
}

type MerchantItemData struct {
	ItemId          string `json:"itemId"`
	Name            string `json:"name"`
	ProductCategory string `json:"productCategory"`
	Price           int    `json:"price"`
	ImageURL        string `json:"imageUrl"`
	CreatedAt       string `json:"createdAt"`
}

type GetMerchantItemsResponse struct {
	Data []MerchantItemData `json:"data"`
	Meta MerchantMeta       `json:"meta"`
}

type NearbyMerchant struct {
	Merchant MerchantData       `json:"merchant"`
	Items    []MerchantItemData `json:"items"`
}

type GetNearbyMerchantsResponse struct {
	Data []NearbyMerchant `json:"data"`
	Meta MerchantMeta     `json:"meta"`
}

func ConvertDomainToDTO(domainMerchant domain.Merchant) MerchantData {
	return MerchantData{
		MerchantID:       domainMerchant.ID,
		Name:             domainMerchant.Name,
		MerchantCategory: string(domainMerchant.MerchantCategory),
		ImageURL:         domainMerchant.ImageURL,
		Location: Location{
			Lat:  domainMerchant.Location.Lat,
			Long: domainMerchant.Location.Long,
		},
		CreatedAt: domainMerchant.CreatedAt,
	}
}

func ConvertDomainItemToDTO(domainItem domain.MerchantItem) MerchantItemData {
	return MerchantItemData{
		ItemId:          domainItem.ID,
		Name:            domainItem.Name,
		ProductCategory: string(domainItem.ProductCategory),
		Price:           int(domainItem.Price),
		ImageURL:        domainItem.ImageURL,
		CreatedAt:       domainItem.CreatedAt,
	}
}
