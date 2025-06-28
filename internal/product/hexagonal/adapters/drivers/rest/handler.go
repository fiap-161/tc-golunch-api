package rest

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/core/model/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/core/ports"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/helper"
	"github.com/gin-gonic/gin"
)

const MaxFileSize = 5 << 20

var allowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

type Handler struct {
	Service ports.ProductService
}

func New(service ports.ProductService) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Create(c *gin.Context) {
	var productDTO dto.ProductRequestDTO

	if err := c.ShouldBindJSON(&productDTO); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	var product model.Product
	product = product.FromRequestDTO(productDTO)
	created, err := h.Service.Create(context.Background(), product)

	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *Handler) ListCategories(c *gin.Context) {
	ctx := context.Background()
	c.JSON(http.StatusOK, h.Service.ListCategories(ctx))
}

func (h *Handler) GetAll(c *gin.Context) {
	query := c.Query("category")
	query = strings.ToUpper(query)
	query = strings.ReplaceAll(query, " ", "")

	ok := enum.IsValidCategory(query)
	if !ok && query != "" {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Validation error",
			MessageError: "Invalid category",
		})
		return
	}

	category := enum.Category(query)
	list, err := h.Service.GetAll(context.Background(), category)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	var products []dto.ProductResponseDTO
	for _, product := range list {
		productDTO := product.FromEntityToResponseDTO()
		products = append(products, productDTO)
	}

	c.JSON(http.StatusOK, dto.ProductListResponseDTO{
		Total: uint(len(list)),
		List:  products,
	})
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var productUpdateDTO dto.ProductRequestUpdateDTO

	if err := c.ShouldBindJSON(&productUpdateDTO); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	var product model.Product
	product = product.FromUpdateDTO(productUpdateDTO)
	updated, err := h.Service.Update(context.Background(), product, id)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.Service.Delete(context.Background(), id)

	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) ValidateIfProductExists(c *gin.Context) {
	id := c.Param("id")

	_, err := h.Service.FindByID(context.Background(), id)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.Next()
}

func (h *Handler) UploadImage(c *gin.Context) {
	uploadDir := os.Getenv("UPLOAD_DIR")
	publicURL := os.Getenv("PUBLIC_URL")

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxFileSize)

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			MessageError: "Validation error",
			Message:      "Image is required or too large (max 5MB).",
		})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, apperror.ErrorDTO{
			MessageError: "Internal Error",
			Message:      "Error opening file.",
		})
		return
	}
	defer file.Close()

	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		c.JSON(http.StatusInternalServerError, apperror.ErrorDTO{
			MessageError: "Internal Error",
			Message:      "Error reading file.",
		})
		return
	}
	contentType := http.DetectContentType(buffer)

	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			MessageError: "Validation error",
			Message:      "Only JPEG and PNG images are allowed.",
		})
		return
	}

	file.Seek(0, 0)

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err = os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, apperror.ErrorDTO{
				MessageError: "Internal Error",
				Message:      "Error creating directory.",
			})
			return
		}
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileHeader.Filename))
	fullPath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(fileHeader, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, apperror.ErrorDTO{
			MessageError: "Internal Error",
			Message:      "Error saving image.",
		})
		return
	}

	imageURL := fmt.Sprintf("%s/uploads/%s", publicURL, filename)
	c.JSON(http.StatusCreated, dto.ImageURLDTO{ImageURL: imageURL})
}
