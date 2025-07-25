package controller

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/cleanarch/external/datasource"
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/cleanarch/gateway"
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/cleanarch/usecases"
)

type Controller struct {
	AdminDatasource datasource.DataSource
}

func Build(productDataSource datasource.DataSource) *Controller {
	return &Controller{
		AdminDatasource: productDataSource}
}

func (c *Controller) Register(ctx context.Context, adminRequest dto.AdminRequestDTO) error {
	adminGateway := gateway.Build(c.AdminDatasource)
	useCase := usecases.Build(*adminGateway)
	admin := dto.FromAdminRequestDTO(adminRequest)
	err := useCase.Create(ctx, admin)

	if err != nil {
		return err
	}

	return nil

}

func (c *Controller) Login(ctx context.Context, adminRequest dto.AdminRequestDTO) (string, bool, error) {
	adminGateway := gateway.Build(c.AdminDatasource)
	useCase := usecases.Build(*adminGateway)
	admin := dto.FromAdminRequestDTO(adminRequest)
	adminId, isAdmin, err := useCase.Login(ctx, admin)

	if err != nil {
		return "", true, err
	}

	return adminId, isAdmin, nil

}
