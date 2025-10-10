package repository

import "errors"

var (
	ErrCalculatedEstimateNotFound      = errors.New("calculated estimate not found")
	ErrCalculatedEstimateAlreadyExists = errors.New("calculated estimate already exists")
	ErrOrderNotFound                   = errors.New("order not found")
	ErrOrderAlreadyExists              = errors.New("order already exists")
	ErrMerchantNotFound                = errors.New("merchant not found")
	ErrMerchantItemNotFound            = errors.New("merchant item not found")
	ErrInvalidUUID                     = errors.New("invalid UUID format")
	ErrInternalServerError             = errors.New("internal server error")
)
