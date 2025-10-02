package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/ProjectSprint-Generalist/BeliMang/internal/db"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	MaxUploadSize = 2 * 1024 * 1024 // 2 MB
	MinUploadSize = 10 * 1024       // 10 KB
)

type ImageHandler struct {
	db  *db.Queries
	min *storage.MinioClient
}

func NewImageHandler(pool *pgxpool.Pool, min *storage.MinioClient) *ImageHandler {
	return &ImageHandler{
		db:  db.New(pool),
		min: min,
	}
}

func (h *ImageHandler) UploadImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize+1024)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid image file. Only .jpg/.jpeg allowed, size must be between 10KB and 2MB",
			Code:    http.StatusBadRequest,
		})
		return
	}
	defer file.Close()

	// Size check
	size := header.Size
	if size < MinUploadSize || size > MaxUploadSize {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid image file. Only .jpg/.jpeg allowed, size must be between 10KB and 2MB",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Extension check
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid image file. Only .jpg/.jpeg allowed, size must be between 10KB and 2MB",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Content-Type check
	head := make([]byte, 512)
	n, _ := file.Read(head)
	contentType := http.DetectContentType(head[:n])
	if contentType != "image/jpeg" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   "Invalid image file. Only .jpg/.jpeg allowed, size must be between 10KB and 2MB",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Reset file pointer
	file.Seek(0, io.SeekStart)

	// Create uuid filename
	id := uuid.New()
	pgUUID := pgtype.UUID{Bytes: id, Valid: true}
	objName := fmt.Sprintf("%s%s", id.String(), ext)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Upload to minio
	uploadURL, err := h.min.PutObject(ctx, objName, file, size, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "failed to upload file",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// Save metadata to db
	image, err := h.db.CreateImage(ctx, db.CreateImageParams{
		ID:        pgUUID,
		Filename:  objName,
		Url:       uploadURL,
		SizeBytes: size,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   "failed to save metadata",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	// response
	c.JSON(http.StatusOK, dto.BaseResponse{
		Message: "File uploaded successfully",
		Data: dto.ImageUploadResponse{
			ImageURL: image.Url,
		},
	})
}
