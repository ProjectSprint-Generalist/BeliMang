package dto

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
}
