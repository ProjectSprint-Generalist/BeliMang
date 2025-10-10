package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/google/uuid"
)

var (
	ErrInvalidOrderData            = errors.New("invalid order data")
	ErrInvalidCalculatedEstimateID = errors.New("invalid calculated estimate ID")
	ErrInvalidUserID               = errors.New("invalid user ID")
)

// Order represents an order in the domain layer
type Order struct {
	ID                   string
	UserID               string
	CalculatedEstimateID string
	CreatedAt            string
}

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

// CalculatedEstimate represents a calculated estimate in the domain layer
type CalculatedEstimate struct {
	ID                           string
	UserID                       string
	TotalPrice                   int32
	EstimatedDeliveryTimeMinutes int32
	EstimateData                 []byte
	CreatedAt                    time.Time
}

// OrderWithEstimate represents an order with its associated estimate data
type OrderWithEstimate struct {
	ID           string
	UserID       string
	EstimateData dto.EstimateRequest
	CreatedAt    time.Time
}

type OrderDetail struct {
	Merchant OrderMerchant `json:"merchant"`
	Items    []OrderItem   `json:"items"`
}

type OrderHistory struct {
	OrderID string
	Orders  []OrderDetail
}

type OrderHistoryResponse []OrderHistory

// NewOrder creates a new order domain entity
func NewOrder(userID, calculatedEstimateID string) (Order, error) {
	if strings.TrimSpace(userID) == "" {
		return Order{}, ErrInvalidUserID
	}

	if strings.TrimSpace(calculatedEstimateID) == "" {
		return Order{}, ErrInvalidCalculatedEstimateID
	}

	// Validate UUID format for calculated estimate ID
	if _, err := uuid.Parse(calculatedEstimateID); err != nil {
		return Order{}, ErrInvalidCalculatedEstimateID
	}

	// Validate UUID format for user ID
	if _, err := uuid.Parse(userID); err != nil {
		return Order{}, ErrInvalidUserID
	}

	return Order{
		UserID:               userID,
		CalculatedEstimateID: calculatedEstimateID,
	}, nil
}

// ValidateOrderCreation validates that an order can be created
func ValidateOrderCreation(userID, calculatedEstimateID string) error {
	_, err := NewOrder(userID, calculatedEstimateID)
	if err != nil {
		return err
	}

	// Additional business validation can be added here
	// For example, checking if the calculated estimate belongs to the user
	// or if the estimate is still valid (not expired)

	return nil
}
