package repository

import "errors"

var (
	ErrMerchantAlreadyExists = errors.New("merchant already exists")
	ErrMerchantNotFound      = errors.New("merchant not found")
	ErrMerchantItemNotFound  = errors.New("merchant item not found")
	ErrInternalServerError   = errors.New("internal server error")
	ErrInvalidUUID           = errors.New("invalid UUID format")
)
