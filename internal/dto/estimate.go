package dto

type EstimateLocation struct {
    Lat  float64 `json:"lat" binding:"required"`
    Long float64 `json:"long" binding:"required"`
}

type EstimateItem struct {
    ItemId   string `json:"itemId" binding:"required,uuid"`
    Quantity int    `json:"quantity" binding:"required,min=1"`
}

type EstimateOrder struct {
    MerchantId      string         `json:"merchantId" binding:"required,uuid"`
    IsStartingPoint bool           `json:"isStartingPoint"`
    Items           []EstimateItem `json:"items" binding:"required,dive"`
}

type EstimateRequest struct {
    UserLocation EstimateLocation `json:"userLocation" binding:"required"`
    Orders       []EstimateOrder  `json:"orders" binding:"required,dive"`
}

type EstimateResponse struct {
    TotalPrice                  float64 `json:"totalPrice"`
    EstimatedDeliveryTimeInMins float64 `json:"estimatedDeliveryTimeInMinutes"`
    CalculatedEstimateID        string  `json:"calculatedEstimateId"`
}