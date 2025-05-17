package rest

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/core/ports"
)

type AdminHandler struct {
	service ports.AdminService
}

func NewAdminHandler(service ports.AdminService) *AdminHandler {
	return &AdminHandler{service: service}
}

func (a *AdminHandler) Register(c *gin.Context) {
	ctx := context.Background()

	var input dto.RegisterDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := a.service.Register(ctx, input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

func (a *AdminHandler) Login(c *gin.Context) {
	ctx := context.Background()

	var input dto.LoginDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.service.Login(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
