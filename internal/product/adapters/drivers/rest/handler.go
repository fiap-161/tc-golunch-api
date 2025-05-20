package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivers/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/ports"
	appError "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/helper"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Service ports.ProductService
}

func NewProductHandler(service ports.ProductService) *ProductHandler {
	return &ProductHandler{Service: service}
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
func (controller *ProductHandler) Create(c *gin.Context) {
	var productDTO dto.ProductRequestDTO

	if err := c.ShouldBindJSON(&productDTO); err != nil {
		c.JSON(http.StatusBadRequest, appError.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	productModel := dto.FromRequestDTOToModel(productDTO)
	product, err := controller.Service.Create(productModel)

	if err != nil {
		helper.HandleError(c, err)
		return
	}

	productRespDTO := dto.FromModelToResponseDTO(product)
	c.JSON(http.StatusCreated, productRespDTO)
}

// List Categories godoc
// @Summary      List Categories
// @Description  List Categories
// @Tags         Product Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Success      200   {array}   enum.CategoryDTO
// @Router       /product/categories [get]
func (controller *ProductHandler) ListCategories(c *gin.Context) {
	c.JSON(http.StatusOK, controller.Service.ListCategories())
}

// Get All Products by Category godoc
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
func (controller *ProductHandler) GetAll(c *gin.Context) {
	query := c.Query("category")
	query = strings.ToLower(query)
	query = strings.ReplaceAll(query, " ", "")

	_, ok := enum.FromCategoryString(query)

	if !ok && query != "" {
		c.JSON(http.StatusBadRequest, appError.ErrorDTO{
			Message:      "Validation error",
			MessageError: "Invalid category",
		})
		return
	}

	list, err := controller.Service.GetAll(query)

	if err != nil {
		helper.HandleError(c, err)
		return
	}

	var products []dto.ProductResponseDTO
	for _, product := range list {
		productRespDTO := dto.FromModelToResponseDTO(product)
		products = append(products, productRespDTO)
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
func (controller *ProductHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	var productUpdateDTO dto.ProductRequestUpdateDTO
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&productUpdateDTO); err != nil {
		c.JSON(http.StatusBadRequest, appError.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	product := dto.FromRequestUpdateDTOToModel(productUpdateDTO)

	productUpdated, err := controller.Service.Update(product, uint(id))

	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.FromModelToResponseDTO(productUpdated))
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
func (controller *ProductHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, _ := strconv.Atoi(idParam)

	err := controller.Service.Delete(uint(id))

	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (controller *ProductHandler) ValidateIfProductExists(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, appError.ErrorDTO{
			Message:      "Validation error",
			MessageError: "ID must be a valid integer",
		})
		return
	}

	_, err2 := controller.Service.FindByID(uint(id))

	if err2 != nil {
		helper.HandleError(c, err2)
		return
	}
	c.Next()
}
