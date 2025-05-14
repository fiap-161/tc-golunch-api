package rest

import (
	"net/http"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivers/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/ports"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/helper"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Service ports.ProductService
}

func NewProductHandler(service ports.ProductService) *ProductHandler {
	return &ProductHandler{Service: service}
}

func (controller *ProductHandler) Create(c *gin.Context) {
	var productDTO dto.ProductRequestDTO
	if err := c.ShouldBindJSON(&productDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":       "Check required fields",
			"message_error": err.Error(),
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
