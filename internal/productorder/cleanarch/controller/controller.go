package controller

import (
	"context"

	orderdto "github.com/fiap-161/tech-challenge-fiap161/internal/order/hexagonal/adapters/drivers/rest/dto"
	productdto "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/external/datasource"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/gateway"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/presenter"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/usecases"
)

type Controller struct {
	productOrderDatasource datasource.DataSource
}

func Build(productDataSource datasource.DataSource) *Controller {
	return &Controller{
		productOrderDatasource: productDataSource}
}

func (c *Controller) CreateBulk(ctx context.Context, listProductOrderRequestDTO []dto.ProductOrderRequestDTO) (int, error) {
	productOrderGateway := gateway.Build(c.productOrderDatasource)
	useCase := usecases.Build(*productOrderGateway)

	var productOrders []entity.ProductOrder
	for _, item := range listProductOrderRequestDTO {
		productOrder := dto.FromRequestDTO(item)
		productOrders = append(productOrders, productOrder)
	}

	length, err := useCase.CreateBulk(ctx, productOrders)
	if err != nil {
		return 0, err
	}

	return length, nil
}

func (c *Controller) FindByOrderID(ctx context.Context, orderId string) ([]dto.ProductOrderResponseDTO, error) {
	productOrderGateway := gateway.Build(c.productOrderDatasource)
	useCase := usecases.Build(*productOrderGateway)
	presenter := presenter.Build()

	productOrderFoundList, findErr := useCase.FindByOrderID(ctx, orderId)
	if findErr != nil {
		return []dto.ProductOrderResponseDTO{}, findErr
	}

	return presenter.FromEntityListToResponseDTOList(productOrderFoundList), nil
}

func (c *Controller) BuildBulkFromOrderAndProducts(
	orderID string,
	orderProductInfo []orderdto.OrderProductInfo,
	productsDTOs []productdto.ProductResponseDTO,
) ([]dto.ProductOrderRequestDTO, error) {
	var result []dto.ProductOrderRequestDTO

	for _, product := range productsDTOs {
		for _, item := range orderProductInfo {
			if product.ID == item.ProductID {
				result = append(result, dto.ProductOrderRequestDTO{
					OrderID:   orderID,
					ProductID: product.ID,
					Quantity:  item.Quantity,
					UnitPrice: product.Price,
				})
			}
		}
	}

	return result, nil
}
