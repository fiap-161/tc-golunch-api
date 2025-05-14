package rest

import (
	"net/http"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivers/dto"
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
// @Accept       json
// @Produce      json
// @Param        request body dto.ProductRequestDTO true "Product to create. Note category is an integer number. See [GET] /product/category"
// @Success      201  {object}  dto.ProductResponseDTO
// @Failure      400  {object}  errors.ErrorDTO
// @Router       /product/ [post]
func (controller *ProductHandler) Create(c *gin.Context) {
	var productDTO dto.ProductRequestDTO
	if err := c.ShouldBindJSON(&productDTO); err != nil {
		c.JSON(http.StatusBadRequest, appError.ErrorDTO{
			Message:      "Check required fields",
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
