package dto

<<<<<<< HEAD
type MerchantLocation struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type MerchantData struct {
	MerchantID       string           `json:"merchantId"`
	Name             string           `json:"name"`
	MerchantCategory string           `json:"merchantCategory"`
	ImageURL         string           `json:"imageUrl"`
	Location         MerchantLocation `json:"location"`
	CreatedAt        string           `json:"createdAt"`
}

type MerchantMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type GetMerchantsResponse struct {
	Data []MerchantData `json:"data"`
	Meta MerchantMeta   `json:"meta"`
=======
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
>>>>>>> a1fde89861b88a9b2953d4083aa2b67c5fee0555
}
