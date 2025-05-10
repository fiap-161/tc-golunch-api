package rest

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/user/core/ports"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(service ports.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) GetUserByID(c *gin.Context) {
	ctx := context.Background()
	id := c.Param("id")

	user, err := u.service.GetUserByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
