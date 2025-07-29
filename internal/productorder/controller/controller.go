package controller

import (
	"context"

	orderdto "github.com/fiap-161/tech-challenge-fiap161/internal/order/hexagonal/adapters/drivers/rest/dto"

	productdto "github.com/fiap-161/tech-challenge-fiap161/internal/product/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/external/datasource"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/gateway"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/presenter"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/usecases"
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
	produtOrderGateway := gateway.Build(c.ProductOrderDatasource)
	useCase := usecases.Build(*produtOrderGateway)
	presenter := presenter.Build()

	productOrderFoundList, err := useCase.FindByOrderID(ctx, orderId)

	if err != nil {
		return []dto.ProductOrderResponseDTO{}, err
	}

	return presenter.FromEntityListToResponseDTOList(productOrderFoundList), nil

}

func (c *Controller) BuildBulkFromOrderAndProducts(
	orderID string,
	orderProductInfo []orderdto.OrderProductInfo,
	productsDTOs []productdto.ProductResponseDTO) ([]dto.ProductOrderRequestDTO, error) {

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
