package dto

// Location - general location type used across the app
type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

// MerchantCategory enum
type MerchantCategory string

const (
	SmallRestaurant       MerchantCategory = "SmallRestaurant"
	MediumRestaurant      MerchantCategory = "MediumRestaurant"
	LargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategory = "BoothKiosk"
	ConvenienceStore      MerchantCategory = "ConvenienceStore"
)

// MerchantCreateRequest for POST /admin/merchants
type MerchantCreateRequest struct {
	Name             string           `json:"name" binding:"required,min=2,max=30"`
	MerchantCategory MerchantCategory `json:"merchantCategory" binding:"required"`
	ImageURL         string           `json:"imageURL" binding:"required"`
	Location         Location         `json:"location" binding:"required"`
}

// MerchantCreateResponse for POST /admin/merchants
type MerchantCreateResponse struct {
	MerchantId string `json:"merchantId"`
}

// MerchantData for GET /admin/merchants response
type MerchantData struct {
	MerchantID       string           `json:"merchantId"`
	Name             string           `json:"name"`
	MerchantCategory string           `json:"merchantCategory"`
	ImageURL         string           `json:"imageUrl"`
	Location         Location         `json:"location"`
	CreatedAt        string           `json:"createdAt"`
}

// MerchantMeta for GET /admin/merchants pagination
type MerchantMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// GetMerchantsResponse for GET /admin/merchants
type GetMerchantsResponse struct {
	Data []MerchantData `json:"data"`
	Meta MerchantMeta   `json:"meta"`
}
