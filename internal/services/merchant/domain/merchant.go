package domain

import (
	"errors"
	"net/url"
	"strings"
)

type MerchantCategory string

const (
	MerchantCategorySmallRestaurant       MerchantCategory = "SmallRestaurant"
	MerchantCategoryMediumRestaurant      MerchantCategory = "MediumRestaurant"
	MerchantCategoryLargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchantCategoryMerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	MerchantCategoryBoothKiosk            MerchantCategory = "BoothKiosk"
	MerchantCategoryConvenienceStore      MerchantCategory = "ConvenienceStore"
)

type ProductCategory string

const (
	ProductCategoryBeverage   ProductCategory = "Beverage"
	ProductCategoryFood       ProductCategory = "Food"
	ProductCategorySnack      ProductCategory = "Snack"
	ProductCategoryCondiments ProductCategory = "Condiments"
	ProductCategoryAdditions  ProductCategory = "Additions"
)

type Location struct {
	Lat  float64
	Long float64
}

var (
	ErrInvalidMerchantCategory = errors.New("invalid merchant category")
	ErrInvalidProductCategory  = errors.New("invalid product category")
	ErrInvalidLocation         = errors.New("invalid location coordinates")
	ErrInvalidName             = errors.New("invalid name")
	ErrInvalidPrice            = errors.New("invalid price")
	ErrInvalidImageURL         = errors.New("invalid image URL")
)

type Merchant struct {
	ID               string
	Name             string
	MerchantCategory MerchantCategory
	Location         Location
	ImageURL         string
	CreatedAt        string
}

type MerchantItem struct {
	ID              string
	MerchantID      string
	Name            string
	ProductCategory ProductCategory
	Price           int32
	ImageURL        string
	CreatedAt       string
}

func NewMerchant(name string, category MerchantCategory, location Location, imageURL string) (Merchant, error) {
	if strings.TrimSpace(name) == "" || len(name) < 2 || len(name) > 30 {
		return Merchant{}, ErrInvalidName
	}

	if !isValidMerchantCategory(category) {
		return Merchant{}, ErrInvalidMerchantCategory
	}

	if err := validateLocation(location); err != nil {
		return Merchant{}, err
	}

	if err := validateImageURL(imageURL); err != nil {
		return Merchant{}, err
	}

	return Merchant{
		Name:             name,
		MerchantCategory: category,
		Location:         location,
		ImageURL:         imageURL,
	}, nil
}

func NewMerchantItem(name string, productCategory ProductCategory, price int32, imageURL string) (MerchantItem, error) {
	if strings.TrimSpace(name) == "" || len(name) < 2 || len(name) > 30 {
		return MerchantItem{}, ErrInvalidName
	}

	if !isValidProductCategory(productCategory) {
		return MerchantItem{}, ErrInvalidProductCategory
	}

	if price < 1 {
		return MerchantItem{}, ErrInvalidPrice
	}

	if err := validateImageURL(imageURL); err != nil {
		return MerchantItem{}, err
	}

	return MerchantItem{
		Name:            name,
		ProductCategory: productCategory,
		Price:           price,
		ImageURL:        imageURL,
	}, nil
}

func isValidMerchantCategory(category MerchantCategory) bool {
	validCategories := map[MerchantCategory]bool{
		MerchantCategorySmallRestaurant:       true,
		MerchantCategoryMediumRestaurant:      true,
		MerchantCategoryLargeRestaurant:       true,
		MerchantCategoryMerchandiseRestaurant: true,
		MerchantCategoryBoothKiosk:            true,
		MerchantCategoryConvenienceStore:      true,
	}

	return validCategories[category]
}

func isValidProductCategory(category ProductCategory) bool {
	validCategories := map[ProductCategory]bool{
		ProductCategoryBeverage:   true,
		ProductCategoryFood:       true,
		ProductCategorySnack:      true,
		ProductCategoryCondiments: true,
		ProductCategoryAdditions:  true,
	}

	return validCategories[category]
}

func validateLocation(location Location) error {
	if location.Lat < -90 || location.Lat > 90 {
		return ErrInvalidLocation
	}
	if location.Long < -180 || location.Long > 180 {
		return ErrInvalidLocation
	}
	return nil
}

func validateImageURL(imageURL string) error {
	if imageURL == "" {
		return ErrInvalidImageURL
	}

	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return ErrInvalidImageURL
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ErrInvalidImageURL
	}

	if parsedURL.Host == "" {
		return ErrInvalidImageURL
	}

	if parsedURL.Path == "" || parsedURL.Path == "/" {
		return ErrInvalidImageURL
	}

	return nil
}
