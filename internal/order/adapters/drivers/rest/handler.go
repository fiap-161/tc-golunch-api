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
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      401  {object}  errors.ErrorDTO
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
	qrCode, err := h.service.Create(ctx, orderDTO)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"qr_code": qrCode,
		"message": "Order created successfully",
	})
}

// Update Order godoc
// @Summary      Update Order
// @Description  Update an existing order status
// @Tags         Order Domain
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Param        request body dto.UpdateOrderDTO true "Order status update"
// @Success      204  "No Content"
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      401  {object}  errors.ErrorDTO
// @Failure      404  {object}  errors.ErrorDTO
// @Router       /order/{id} [put]
func (h *handler) Update(c *gin.Context) {
	id := c.Param("id")

	var orderUpdate dto.UpdateOrderDTO
	if err := c.ShouldBindJSON(&orderUpdate); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	err := h.service.Update(context.Background(), id, orderUpdate.Status)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetAll Get All Orders godoc
// @Summary      Get all orders
// @Description  Get all orders
// @Tags         Order Domain
// @Security 	 BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      401  {object}  errors.ErrorDTO
// @Router       /order/ [get]
func (h *handler) GetAll(c *gin.Context) {
	ctx := context.Background()
	orders, err := h.service.GetAll(ctx)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

// GetPanel Get Order Panel godoc
// @Summary      Get Order Panel
// @Description  Get the order panel with all orders that are in the panel status
// @Tags         Order Domain
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.OrderPanelDTO
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      401  {object}  errors.ErrorDTO
// @Router       /order/panel [get]
func (h *handler) GetPanel(c *gin.Context) {
	ctx := context.Background()
	orders, err := h.service.GetPanel(ctx)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	var panel []dto.OrderPanelItemDTO
	for _, order := range orders {
		panelDTO := order.ToPanelItemDTO()
		panel = append(panel, panelDTO)
	}

	c.JSON(http.StatusOK, dto.OrderPanelDTO{
		Orders: panel,
	})
}
