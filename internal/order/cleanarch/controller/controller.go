package controller

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/external/datasource"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/gateway"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/presenter"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/usecases"
)

type ProductService interface {
	FindByIDs(ctx context.Context, ids []string) ([]dto.ProductDTO, error)
}
type ProductOrderService interface {
	BuildBulkFromOrderAndProducts(orderID string, items []dto.OrderProductInfo, products []dto.ProductDTO) ([]dto.ProductOrderDTO, error)
	CreateBulk(ctx context.Context, productOrders []dto.ProductOrderDTO) ([]dto.ProductOrderDTO, error)
}
type PaymentService interface {
	CreateByOrderID(ctx context.Context, orderID string) (dto.PaymentDTO, error)
}

type Controller struct {
	OrderDatasource     datasource.DataSource
	ProductService      ProductService
	ProductOrderService ProductOrderService
	PaymentService      PaymentService
}

func Build(orderDatasource datasource.DataSource, productService ProductService, productOrderService ProductOrderService, paymentService PaymentService) *Controller {
	return &Controller{
		OrderDatasource:     orderDatasource,
		ProductService:      productService,
		ProductOrderService: productOrderService,
		PaymentService:      paymentService,
	}
}

func (c *Controller) Create(ctx context.Context, orderDTO dto.CreateOrderDTO) (string, error) {
	orderGateway := gateway.Build(c.OrderDatasource)
	useCase := usecases.Build(orderGateway, c.ProductService, c.ProductOrderService, c.PaymentService)

	return useCase.CreateCompleteOrder(ctx, orderDTO)
}

func (c *Controller) GetAll(ctx context.Context) ([]dto.OrderDAO, error) {
	orderGateway := gateway.Build(c.OrderDatasource)
	useCase := usecases.Build(orderGateway, c.ProductService, c.ProductOrderService, c.PaymentService)
	presenter := presenter.Build()

	orders, err := useCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return presenter.FromEntityListToDAOList(orders), nil
}

func (c *Controller) GetPanel(ctx context.Context, status []string) ([]dto.OrderDAO, error) {
	orderGateway := gateway.Build(c.OrderDatasource)
	useCase := usecases.Build(orderGateway, c.ProductService, c.ProductOrderService, c.PaymentService)
	presenter := presenter.Build()

	orders, err := useCase.GetPanel(ctx, status)
	if err != nil {
		return nil, err
	}
	return presenter.FromEntityListToDAOList(orders), nil
}

func (c *Controller) FindByID(ctx context.Context, id string) (dto.OrderDAO, error) {
	orderGateway := gateway.Build(c.OrderDatasource)
	useCase := usecases.Build(orderGateway, c.ProductService, c.ProductOrderService, c.PaymentService)
	presenter := presenter.Build()

	order, err := useCase.FindByID(ctx, id)
	if err != nil {
		return dto.OrderDAO{}, err
	}
	return presenter.FromEntityToDAO(order), nil
}

func (c *Controller) Update(ctx context.Context, orderDTO dto.OrderDAO) (dto.OrderDAO, error) {
	orderGateway := gateway.Build(c.OrderDatasource)
	useCase := usecases.Build(orderGateway, c.ProductService, c.ProductOrderService, c.PaymentService)
	presenter := presenter.Build()

	order := dto.FromOrderDAO(orderDTO)
	updated, err := useCase.Update(ctx, order)
	if err != nil {
		return dto.OrderDAO{}, err
	}
	return presenter.FromEntityToDAO(updated), nil
}
