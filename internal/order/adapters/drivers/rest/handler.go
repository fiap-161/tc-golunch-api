package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/core/ports"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/helper"
)

type handler struct {
	service ports.OrderService
}

func New(service ports.OrderService) *handler {
	return &handler{
		service: service,
	}
}

// Create Create Order godoc
// @Summary      Create Order
// @Description  Create a new order
// @Tags         Order Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateOrderDTO true "Order to create. Note that the customer_id is automatically set from the authenticated user."
// @Success      200  {object}  gin.H{"id": "string", "message": "Order created successfully"}
// @Failure      400  {object}  apperror.ErrorDTO
// @Failure      401  {object}  apperror.ErrorDTO
// @Router       /order/ [post]
func (h *handler) Create(c *gin.Context) {
	ctx := context.Background()

	var orderDTO dto.CreateOrderDTO
	if err := c.ShouldBindJSON(&orderDTO); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	validateErr := orderDTO.Validate()
	if validateErr != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "validation failed",
			MessageError: validateErr.Error(),
		})
		return
	}

	customerIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, apperror.ErrorDTO{
			Message:      "unauthorized",
			MessageError: "user id not found in context",
		})
		return
	}
	customerID := customerIDRaw.(string)

	orderDTO.CustomerID = customerID
	id, err := h.service.Create(ctx, orderDTO)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "Order created successfully",
	})
}

// TODO Create Update Method for kitchen handler ready and completed orders

// GetAll Get All Orders godoc
// @Summary      Get all orders
// @Description  Get all orders
// @Tags         Order Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  gin.H{"orders": []dto.OrderResponseDTO}
// @Failure      400  {object}  apperror.ErrorDTO
// @Failure      401  {object}  apperror.ErrorDTO
// @Router       /order/ [get]
func (h *handler) GetAll(c *gin.Context) {
	ctx := context.Background()
	products, err := h.service.GetAll(ctx)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}
