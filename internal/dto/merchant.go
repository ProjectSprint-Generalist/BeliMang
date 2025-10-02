package dto

type MerchantCategory string

const (
	SmallRestaurant       MerchantCategory = "SmallRestaurant"
	MediumRestaurant      MerchantCategory = "MediumRestaurant"
	LargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategory = "BoothKiosk"
	ConvenienceStore      MerchantCategory = "ConvenienceStore"
)

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

type MerchantCreateRequest struct {
	Name             string           `json:"name" binding:"required,min=2,max=30"`
	MerchantCategory MerchantCategory `json:"merchantCategory" binding:"required"`
	ImageURL         string           `json:"imageURL" binding:"required"`
	Location         Location         `json:"location" binding:"required"`
}

type MerchantCreateResponse struct {
	MerchantId string `json:"merchantId"`
}
