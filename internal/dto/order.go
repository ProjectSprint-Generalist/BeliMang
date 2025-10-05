package dto

// Request DTOs
type CreateOrderRequest struct {
	CalculatedEstimateID string `json:"calculatedEstimateId" binding:"required,uuid"`
}

// Response DTOs
type CreateOrderResponse struct {
	OrderID string `json:"orderId"`
}

// Order history DTOs
type OrderLocation struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type OrderMerchant struct {
	MerchantID       string        `json:"merchantId"`
	Name             string        `json:"name"`
	MerchantCategory string        `json:"merchantCategory"`
	ImageURL         string        `json:"imageUrl"`
	Location         OrderLocation `json:"location"`
	CreatedAt        string        `json:"createdAt"`
}

type OrderItem struct {
	ItemID          string `json:"itemId"`
	Name            string `json:"name"`
	ProductCategory string `json:"productCategory"`
	Price           int    `json:"price"`
	Quantity        int    `json:"quantity"`
	ImageURL        string `json:"imageUrl"`
	CreatedAt       string `json:"createdAt"`
}

type OrderDetail struct {
	Merchant OrderMerchant `json:"merchant"`
	Items    []OrderItem   `json:"items"`
}

type OrderHistory struct {
	OrderID string        `json:"orderId"`
	Orders  []OrderDetail `json:"orders"`
}

// Response with meta pagination info
type GetOrdersResponseMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type GetOrdersResponse struct {
	Data []OrderHistory        `json:"data"`
	Meta GetOrdersResponseMeta `json:"meta"`
}

// Query parameters for GET orders
type GetOrdersParams struct {
	MerchantID       *string `form:"merchantId"`
	Limit            int     `form:"limit"`
	Offset           int     `form:"offset"`
	Name             *string `form:"name"`
	MerchantCategory *string `form:"merchantCategory"`
	CreatedAt        *string `form:"createdAt"`
}
