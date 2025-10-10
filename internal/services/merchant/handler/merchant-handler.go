package handler

import (
	"net/http"
	"net/url"
	"strconv"

	base "github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/merchant/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/merchant/service"
	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	svc *service.MerchantService
}

func NewMerchantHandler(svc *service.MerchantService) *MerchantHandler {
	return &MerchantHandler{svc: svc}
}

func isValidImageURL(imageURL string) bool {
	if imageURL == "" {
		return false
	}

	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return false
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}

	if parsedURL.Host == "" {
		return false
	}

	if parsedURL.Path == "" || parsedURL.Path == "/" {
		return false
	}

	return true
}

func (h *MerchantHandler) CreateMerchant(c *gin.Context) {
	var payload dto.MerchantCreateRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, base.ErrorResponse{
			Success: false,
			Error:   "Invalid input: please make sure you have provided a valid name, merchant category, image URL, and location",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if !dto.ValidMerchantCategories[payload.MerchantCategory] {
		c.JSON(http.StatusBadRequest, base.ErrorResponse{
			Success: false,
			Error:   "Invalid merchant category. Must be one of: SmallRestaurant, MediumRestaurant, LargeRestaurant, MerchandiseRestaurant, BoothKiosk, ConvenienceStore",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if payload.Location.Lat < -90 || payload.Location.Lat > 90 || payload.Location.Lat == 0 {
		c.JSON(http.StatusBadRequest, base.ErrorResponse{
			Success: false,
			Error:   "Invalid latitude. Must be between -90 and 90, and not zero",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if payload.Location.Long < -180 || payload.Location.Long > 180 || payload.Location.Long == 0 {
		c.JSON(http.StatusBadRequest, base.ErrorResponse{
			Success: false,
			Error:   "Invalid longitude. Must be between -180 and 180, and not zero",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if !isValidImageURL(payload.ImageURL) {
		c.JSON(http.StatusBadRequest, base.ErrorResponse{
			Success: false,
			Error:   "Invalid image URL. Must be a complete HTTP/HTTPS URL with a path (e.g., https://example.com/image.jpg)",
			Code:    http.StatusBadRequest,
		})
		return
	}

	resp, err := h.svc.CreateMerchant(c.Request.Context(), payload)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrInvalidMerchantData || err == service.ErrMerchantNotFound {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, base.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    statusCode,
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *MerchantHandler) GetMerchants(c *gin.Context) {
	merchantID := c.Query("merchantId")
	name := c.Query("name")
	merchantCategory := c.Query("merchantCategory")
	createdAt := c.Query("createdAt")

	limit := int32(5)
	if limitStr := c.Query("limit"); limitStr != "" {
		if val, err := strconv.ParseInt(limitStr, 10, 32); err == nil && val > 0 {
			limit = int32(val)
		}
	}

	offset := int32(0)
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if val, err := strconv.ParseInt(offsetStr, 10, 32); err == nil && val >= 0 {
			offset = int32(val)
		}
	}

	resp, err := h.svc.GetMerchants(c.Request.Context(), merchantID, name, merchantCategory, createdAt, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, base.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MerchantHandler) CreateMerchantItem(c *gin.Context) {
	merchantID := c.Param("merchantId")

	var payload dto.MerchantItemCreateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, base.ErrorResponse{
			Success: false,
			Error:   "Invalid input: please make sure you have provided valid name, product category, price, and image URL",
			Code:    http.StatusBadRequest,
		})
		return
	}

	resp, err := h.svc.CreateMerchantItem(c.Request.Context(), merchantID, payload)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrMerchantNotFound || err == service.ErrInvalidMerchantItemData {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, base.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    statusCode,
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *MerchantHandler) GetMerchantItems(c *gin.Context) {
	merchantID := c.Param("merchantId")

	itemID := c.Query("itemId")
	name := c.Query("name")
	productCategory := c.Query("productCategory")
	createdAt := c.Query("createdAt")

	limit := int32(5)
	if limitStr := c.Query("limit"); limitStr != "" {
		if val, err := strconv.ParseInt(limitStr, 10, 32); err == nil && val > 0 {
			limit = int32(val)
		}
	}

	offset := int32(0)
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if val, err := strconv.ParseInt(offsetStr, 10, 32); err == nil && val >= 0 {
			offset = int32(val)
		}
	}

	resp, err := h.svc.GetMerchantItems(c.Request.Context(), merchantID, itemID, name, productCategory, createdAt, limit, offset)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == service.ErrMerchantNotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, base.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    statusCode,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MerchantHandler) GetNearbyMerchants(c *gin.Context) {
	coords := c.Param("coords")
	var lat, long float64

	if coords != "" {
		commaIdx := -1
		for i := 0; i < len(coords); i++ {
			if coords[i] == ',' {
				commaIdx = i
				break
			}
		}
		if commaIdx <= 0 || commaIdx >= len(coords)-1 {
			c.JSON(http.StatusBadRequest, base.ErrorResponse{
				Success: false,
				Error:   "lat/long is not valid",
				Code:    http.StatusBadRequest,
			})
			return
		}

		latStr := coords[:commaIdx]
		longStr := coords[commaIdx+1:]

		latParsed, err1 := strconv.ParseFloat(latStr, 64)
		longParsed, err2 := strconv.ParseFloat(longStr, 64)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, base.ErrorResponse{
				Success: false,
				Error:   "lat/long is not valid",
				Code:    http.StatusBadRequest,
			})
			return
		}

		if latParsed < -90 || latParsed > 90 || longParsed < -180 || longParsed > 180 {
			c.JSON(http.StatusBadRequest, base.ErrorResponse{
				Success: false,
				Error:   "lat/long is not valid",
				Code:    http.StatusBadRequest,
			})
			return
		}

		lat = latParsed
		long = longParsed
	} else {
		c.JSON(http.StatusBadRequest, base.ErrorResponse{
			Success: false,
			Error:   "lat/long is not valid",
			Code:    http.StatusBadRequest,
		})
		return
	}

	merchantId := c.Query("merchantId")
	name := c.Query("name")
	merchantCategory := c.Query("merchantCategory")

	limit := int32(5)
	if limitStr := c.Query("limit"); limitStr != "" {
		if val, err := strconv.ParseInt(limitStr, 10, 32); err == nil && val > 0 {
			limit = int32(val)
		}
	}

	offset := int32(0)
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if val, err := strconv.ParseInt(offsetStr, 10, 32); err == nil && val >= 0 {
			offset = int32(val)
		}
	}

	resp, err := h.svc.GetNearbyMerchants(c.Request.Context(), lat, long, merchantId, name, merchantCategory, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, base.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
