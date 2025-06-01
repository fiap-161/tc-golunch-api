package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/core/ports"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type handler struct {
	service ports.PaymentService
}

func New(service ports.PaymentService) *handler {
	return &handler{
		service: service,
	}
}

// CheckPayment godoc
// @Summary      Check Payment [Mercado Pago Integration]
// @Description  Check the status of a payment by its resource URL
// @Tags         Payment Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.CheckPaymentDTO true "Resource URL to check payment status"
// @Success      200
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      500  {object}  errors.ErrorDTO
// @Router       /payment/check [post]
func (h *handler) CheckPayment(c *gin.Context) {
	ctx := c.Request.Context()

	var checkPaymentDTO dto.CheckPaymentDTO
	if err := c.ShouldBindJSON(&checkPaymentDTO); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	if err := c.ShouldBindJSON(&checkPaymentDTO); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "invalid request body",
			MessageError: err.Error(),
		})
	}

	_, err := h.service.CheckPayment(ctx, checkPaymentDTO.Resource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apperror.ErrorDTO{
			Message:      "failed to verify payment",
			MessageError: err.Error(),
		})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
