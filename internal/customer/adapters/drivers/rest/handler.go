package rest

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/ports"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CustomerHandler struct {
	service ports.CustomerService
}

func NewCustomerHandler(service ports.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) Create(c *gin.Context) {
	ctx := context.Background()

	var customerDTO dto.CreateCustomerDTO
	if err := c.ShouldBindJSON(&customerDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.Create(ctx, customerDTO)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "Customer created successfully",
	})
}

func (ch *CustomerHandler) Identify(c *gin.Context) {
	ctx := context.Background()
	CPF := c.Param("cpf")

	token, err := ch.service.Identify(ctx, CPF)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}

func (ch *CustomerHandler) Anonymous(c *gin.Context) {
	ctx := context.Background()

	token, err := ch.service.Identify(ctx, "")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}
