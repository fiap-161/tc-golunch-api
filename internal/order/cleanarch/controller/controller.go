package controller

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/external/datasource"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/gateway"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/presenter"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/usecases"
	paymentuc "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/usecases"
	productuc "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/usecases"
	productorderuc "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/usecases"
)

type Controller struct {
	orderDatasource     datasource.DataSource
	productUseCase      productuc.UseCases
	productOrderUseCase productorderuc.UseCases
	paymentUseCase      paymentuc.UseCases
}

func Build(
	orderDatasource datasource.DataSource,
	productService productuc.UseCases,
	productOrderService productorderuc.UseCases,
	paymentService paymentuc.UseCases,
) *Controller {
	return &Controller{
		orderDatasource:     orderDatasource,
		productUseCase:      productService,
		productOrderUseCase: productOrderService,
		paymentUseCase:      paymentService,
	}
}

func (c *Controller) Create(ctx context.Context, orderDTO dto.CreateOrderDTO) (string, error) {
	orderGateway := gateway.Build(c.orderDatasource)
	useCase := usecases.Build(orderGateway, c.productUseCase, c.productOrderUseCase, c.paymentUseCase)

	return useCase.CreateCompleteOrder(ctx, orderDTO)
}

func (c *Controller) GetAll(ctx context.Context) ([]dto.OrderDAO, error) {
	orderGateway := gateway.Build(c.orderDatasource)
	useCase := usecases.Build(orderGateway, c.productUseCase, c.productOrderUseCase, c.paymentUseCase)
	presenter := presenter.Build()

	orders, err := useCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return presenter.FromEntityListToDAOList(orders), nil
}

func (c *Controller) GetPanel(ctx context.Context, status []string) ([]dto.OrderDAO, error) {
	orderGateway := gateway.Build(c.orderDatasource)
	useCase := usecases.Build(orderGateway, c.productUseCase, c.productOrderUseCase, c.paymentUseCase)
	presenter := presenter.Build()

	orders, err := useCase.GetPanel(ctx, status)
	if err != nil {
		return nil, err
	}

	return presenter.FromEntityListToDAOList(orders), nil
}

func (c *Controller) FindByID(ctx context.Context, id string) (dto.OrderDAO, error) {
	orderGateway := gateway.Build(c.orderDatasource)
	useCase := usecases.Build(orderGateway, c.productUseCase, c.productOrderUseCase, c.paymentUseCase)
	presenter := presenter.Build()

	order, err := useCase.FindByID(ctx, id)
	if err != nil {
		return dto.OrderDAO{}, err
	}

	return presenter.FromEntityToDAO(order), nil
}

func (c *Controller) Update(ctx context.Context, orderDTO dto.OrderDAO) (dto.OrderDAO, error) {
	orderGateway := gateway.Build(c.orderDatasource)
	useCase := usecases.Build(orderGateway, c.productUseCase, c.productOrderUseCase, c.paymentUseCase)
	presenter := presenter.Build()

	order := dto.FromOrderDAO(orderDTO)
	updated, err := useCase.Update(ctx, order)
	if err != nil {
		return dto.OrderDAO{}, err
	}

	return presenter.FromEntityToDAO(updated), nil
}
