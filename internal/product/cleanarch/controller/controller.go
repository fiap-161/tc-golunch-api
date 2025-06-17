package controller

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/external/datasource"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/gateway"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/presenter"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/usecases"
)

type Controller struct {
	ProductDatasource datasource.DataSource
	Presenter         presenter.Presenter
}

func Build(productDataSource datasource.DataSource, presenter presenter.Presenter) *Controller {
	return &Controller{
		ProductDatasource: productDataSource,
		Presenter:         presenter,
	}
}

func (c *Controller) Create(ctx context.Context, productDTO dto.ProductRequestDTO) (dto.ProductResponseDTO, error) {
	productGateway := gateway.Build(c.ProductDatasource)
	useCase := usecases.Build(*productGateway)
	product, _ := useCase.CreateProduct(ctx, productDTO)
	return c.Presenter.FromEntityToResponseDTO(product), nil
}
