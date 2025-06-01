package rest

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
	"net/http"
	"strings"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/ports"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/helper"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service ports.ProductService
}

func New(service ports.ProductService) *Handler {
	return &Handler{Service: service}
}

// Create Product godoc
// @Summary      Create Product
// @Description  Create a new product
// @Tags         Product Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.ProductRequestDTO true "Product to create. Note category is an integer number. See [GET] /product/categories to get a valid category_id"
// @Success      201  {object}  dto.ProductResponseDTO
// @Failure      400  {object}  errors.ErrorDTO
// @Router       /product/ [post]
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

// ListCategories List Categories godoc
// @Summary      List Categories
// @Description  List Categories
// @Tags         Product Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Success      200   {array}   enum.CategoryDTO
// @Router       /product/categories [get]
func (h *Handler) ListCategories(c *gin.Context) {
	ctx := context.Background()
	c.JSON(http.StatusOK, h.Service.ListCategories(ctx))
}

// GetAll Get All Products by Category godoc
// @Summary      Get all products by category
// @Description  Returns all products. Optionally, filter by category using query param. Categories must match those returned from [GET] /product/categories.
// @Tags         Product Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        category query string false "Category name (e.g., 'bebida', 'lanche', 'acompanhamento', 'sobremesa')"
// @Success      200  {object}  dto.ProductListResponseDTO
// @Failure      400  {object}  errors.ErrorDTO
// @Router       /product [get]
func (h *Handler) GetAll(c *gin.Context) {
	query := c.Query("category")
	query = strings.ToLower(query)
	query = strings.ReplaceAll(query, " ", "")

	_, ok := enum.FromCategoryString(query)

	if !ok && query != "" {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Validation error",
			MessageError: "Invalid category",
		})
		return
	}

	list, err := h.Service.GetAll(context.Background(), query)
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

// Update Product godoc
// @Summary      Update Product
// @Description  Update an existing product by ID
// @Tags         Product Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      int                         true  "Product ID"
// @Param        request  body      dto.ProductRequestUpdateDTO true  "Product data to update"
// @Success      200      {object}  dto.ProductResponseDTO
// @Failure      400      {object}  errors.ErrorDTO
// @Router       /product/{id} [put]
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

// Delete Product godoc
// @Summary      Delete Product
// @Description  Delete a product by ID
// @Tags         Product Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      204  "No Content"
// @Failure      400  {object}  errors.ErrorDTO
// @Router       /product/{id} [delete]
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
