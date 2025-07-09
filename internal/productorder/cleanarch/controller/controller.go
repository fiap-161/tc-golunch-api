package controller

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/external/datasource"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/gateway"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/usecases"
)

type Controller struct {
	ProductOrderDatasource datasource.DataSource
}

func Build(productDataSource datasource.DataSource) *Controller {
	return &Controller{
		ProductOrderDatasource: productDataSource}
}

func (c *Controller) CreateBulk(ctx context.Context, listProductOrderRequestDTO []dto.ProductOrderRequestDTO) (int, error) {
	produtOrderGateway := gateway.Build(c.ProductOrderDatasource)
	useCase := usecases.Build(*produtOrderGateway)

	var productOrders []entity.ProductOrder
	for _, item := range listProductOrderRequestDTO {
		productOrder := entity.FromRequestDTO(item)
		productOrders = append(productOrders, productOrder)
	}

	length, err := useCase.CreateBulk(ctx, productOrders)

	if err != nil {
		return 0, err
	}

	return length, nil

}
