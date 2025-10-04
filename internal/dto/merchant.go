package dto

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type MerchantCategory string

const (
	SmallRestaurant       MerchantCategory = "SmallRestaurant"
	MediumRestaurant      MerchantCategory = "MediumRestaurant"
	LargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategory = "BoothKiosk"
	ConvenienceStore      MerchantCategory = "ConvenienceStore"
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

// ProductCategory enum
type ProductCategory string

const (
	Beverage  ProductCategory = "Beverage"
	Food      ProductCategory = "Food"
	Snack     ProductCategory = "Snack"
	Condiment ProductCategory = "Condiment"
	Additions ProductCategory = "Additions"
)

// MerchantItemCreateRequest for POST /admin/merchants/:merchantId/items
type MerchantItemCreateRequest struct {
	Name            string          `json:"name" binding:"required,min=2,max=30"`
	ProductCategory ProductCategory `json:"productCategory" binding:"required"`
	Price           int             `json:"price" binding:"required,min=1"`
	ImageURL        string          `json:"imageUrl" binding:"required"`
}

// MerchantItemCreateResponse for POST /admin/merchants/:merchantId/items
type MerchantItemCreateResponse struct {
	ItemId string `json:"itemId"`
}

// MerchantItemData for GET /admin/merchants/:merchantId/items response
type MerchantItemData struct {
	ItemId          string `json:"itemId"`
	Name            string `json:"name"`
	ProductCategory string `json:"productCategory"`
	Price           int    `json:"price"`
	ImageURL        string `json:"imageUrl"`
	CreatedAt       string `json:"createdAt"`
}

// GetMerchantItemsResponse for GET /admin/merchants/:merchantId/items
type GetMerchantItemsResponse struct {
	Data []MerchantItemData `json:"data"`
	Meta MerchantMeta       `json:"meta"`
}

type GetNearbyMerchantsRequest struct {
	MerchantId       string           `json:"merchantId"`
	Name             string           `json:"name"`
	MerchantCategory MerchantCategory `json:"merchantCategory"`
}
