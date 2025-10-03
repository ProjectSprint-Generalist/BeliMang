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
