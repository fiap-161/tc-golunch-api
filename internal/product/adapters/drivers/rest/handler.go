package rest

import (
	"net/http"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivers/dto"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
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
	c.JSON(http.StatusCreated, productDTO)
}
