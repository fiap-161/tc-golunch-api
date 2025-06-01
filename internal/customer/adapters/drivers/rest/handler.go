package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/ports"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/helper"
)

type CustomerHandler struct {
	service ports.CustomerService
}

func NewCustomerHandler(service ports.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

// Create godoc
// @Summary      Cria um novo cliente
// @Description  Cria um cliente com base nas informações enviadas no corpo da requisição
// @Tags         Customer Domain
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateCustomerDTO  true  "Dados do cliente"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  errors.ErrorDTO
// @Failure      500      {object}  errors.ErrorDTO
// @Router       /customer/register [post]
func (h *CustomerHandler) Create(c *gin.Context) {
	ctx := context.Background()

	var customerDTO dto.CreateCustomerDTO
	if err := c.ShouldBindJSON(&customerDTO); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	id, err := h.service.Create(ctx, customerDTO)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "Customer created successfully",
	})
}

// Identify godoc
// @Summary      Identifica cliente por CPF
// @Description  Retorna um token JWT ao identificar o cliente pelo CPF
// @Tags         Customer Domain
// @Accept       json
// @Produce      json
// @Param        cpf   path      string     true  "CPF do cliente"
// @Success      200   {object}  TokenDTO
// @Failure      404   {object}  errors.ErrorDTO
// @Failure      500   {object}  errors.ErrorDTO
// @Router       /customer/identify/{cpf} [get]
func (h *CustomerHandler) Identify(c *gin.Context) {
	ctx := context.Background()
	CPF := c.Param("cpf")

	token, err := h.service.Identify(ctx, CPF)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, &TokenDTO{
		TokenString: token,
	})
}

// Anonymous godoc
// @Summary      Gera cliente anônimo
// @Description  Gera um token JWT para um cliente anônimo (sem CPF)
// @Tags         Customer Domain
// @Accept       json
// @Produce      json
// @Success      200  {object}  TokenDTO
// @Failure      500  {object}  errors.ErrorDTO
// @Router       /customer/anonymous [get]
func (h *CustomerHandler) Anonymous(c *gin.Context) {
	ctx := context.Background()

	token, err := h.service.Identify(ctx, "")
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, &TokenDTO{
		TokenString: token,
	})
}

type TokenDTO struct {
	TokenString string `json:"token"`
}
