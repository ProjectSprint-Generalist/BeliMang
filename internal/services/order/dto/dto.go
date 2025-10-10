package dto

type Location struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

// UserLocation represents the user's location for estimates
type UserLocation struct {
	Lat  float64 `json:"lat" binding:"required"`
	Long float64 `json:"long" binding:"required"`
}

// ItemRequest represents an item in an order request
type ItemRequest struct {
	ItemID   string `json:"itemId" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
}

// OrderRequest represents a merchant order in an estimate request
type OrderRequest struct {
	MerchantID      string        `json:"merchantId" binding:"required"`
	IsStartingPoint bool          `json:"isStartingPoint"`
	Items           []ItemRequest `json:"items" binding:"required,dive"`
}

// EstimateRequest represents the request for getting a price estimate
type EstimateRequest struct {
	UserLocation Location       `json:"userLocation" binding:"required"`
	Orders       []OrderRequest `json:"orders" binding:"required,min=1,dive"`
}

// EstimateResponse represents the response from estimate endpoint
type EstimateResponse struct {
	TotalPrice                     int    `json:"totalPrice"`
	EstimatedDeliveryTimeInMinutes int    `json:"estimatedDeliveryTimeInMinutes"`
	CalculatedEstimateID           string `json:"calculatedEstimateId"`
}

// CreateOrderRequest represents the request for placing an order
type CreateOrderRequest struct {
	CalculatedEstimateID string `json:"calculatedEstimateId" binding:"required"`
}

// CreateOrderResponse represents the response from order creation
type CreateOrderResponse struct {
	OrderID string `json:"orderId"`
}

// OrderLocation represents location in order history response
type OrderLocation struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

// OrderMerchant represents merchant details in order history
type OrderMerchant struct {
	MerchantID       string        `json:"merchantId"`
	Name             string        `json:"name"`
	MerchantCategory string        `json:"merchantCategory"`
	ImageURL         string        `json:"imageUrl"`
	Location         OrderLocation `json:"location"`
	CreatedAt        string        `json:"createdAt"`
}

// OrderItem represents item details in order history
type OrderItem struct {
	ItemID          string `json:"itemId"`
	Name            string `json:"name"`
	ProductCategory string `json:"productCategory"`
	Price           int    `json:"price"`
	Quantity        int    `json:"quantity"`
	ImageURL        string `json:"imageUrl"`
	CreatedAt       string `json:"createdAt"`
}

// OrderDetail represents detailed order information
type OrderDetail struct {
	Merchant OrderMerchant `json:"merchant"`
	Items    []OrderItem   `json:"items"`
}

// OrderHistory represents a single order in history
type OrderHistory struct {
	OrderID string        `json:"orderId"`
	Orders  []OrderDetail `json:"orders"`
}

// OrderHistoryResponse represents the complete order history response
type OrderHistoryResponse struct {
	Data []OrderHistory `json:"data"`
	Meta OrderMeta      `json:"meta"`
}

// OrderMeta represents pagination metadata for order history
type OrderMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// GetOrdersParams represents query parameters for getting orders
type GetOrdersParams struct {
	MerchantID       *string `form:"merchantId"`
	Limit            int     `form:"limit"`
	Offset           int     `form:"offset"`
	Name             *string `form:"name"`
	MerchantCategory *string `form:"merchantCategory"`
	CreatedAt        *string `form:"createdAt"`
}

// Merchant categories enum for filtering
type MerchantCategory string

const (
	MerchantCategorySmallRestaurant       MerchantCategory = "SmallRestaurant"
	MerchantCategoryMediumRestaurant      MerchantCategory = "MediumRestaurant"
	MerchantCategoryLargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchantCategoryMerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	MerchantCategoryBoothKiosk            MerchantCategory = "BoothKiosk"
	MerchantCategoryConvenienceStore      MerchantCategory = "ConvenienceStore"
)

// ValidMerchantCategories contains all valid merchant category values
var ValidMerchantCategories = map[MerchantCategory]bool{
	MerchantCategorySmallRestaurant:       true,
	MerchantCategoryMediumRestaurant:      true,
	MerchantCategoryLargeRestaurant:       true,
	MerchantCategoryMerchandiseRestaurant: true,
	MerchantCategoryBoothKiosk:            true,
	MerchantCategoryConvenienceStore:      true,
}
