package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/core/ports"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/helper"
)

type AdminHandler struct {
	service ports.AdminService
}

func NewAdminHandler(service ports.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

// Register godoc
// @Summary      Register Admin
// @Description  Register a new admin user
// @Tags         Admin Domain
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RegisterDTO  true  "Admin registration details"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  errors.ErrorDTO
// @Failure      500      {object}  errors.ErrorDTO
// @Router       /admin/register [post]
func (a *AdminHandler) Register(c *gin.Context) {
	ctx := context.Background()

	var input dto.RegisterDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	err := a.service.Register(ctx, input)

	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// Login godoc
// @Summary      Admin Login
// @Description  Authenticates an admin user and returns a JWT token
// @Tags         Admin Domain
// @Accept       json
// @Produce      json
// @Param        request  body      dto.LoginDTO  true  "Admin login credentials"
// @Success      200      {object}  TokenDTO
// @Failure      400      {object}  errors.ErrorDTO
// @Failure      401      {object}  errors.ErrorDTO
// @Failure      500      {object}  errors.ErrorDTO
// @Router       /admin/login [post]
func (a *AdminHandler) Login(c *gin.Context) {
	ctx := context.Background()

	var input dto.LoginDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	token, err := a.service.Login(ctx, input)
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
