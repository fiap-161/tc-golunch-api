package controller

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/external/datasource"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/gateway"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/presenter"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/usecases"
)

type ProductService interface {
	FindByIDs(ctx context.Context, ids []string) ([]ProductDTO, error)
}
type ProductOrderService interface {
	BuildBulkFromOrderAndProducts(orderID string, items []dto.OrderProductInfo, products []ProductDTO) ([]ProductOrderDTO, error)
	CreateBulk(ctx context.Context, productOrders []ProductOrderDTO) ([]ProductOrderDTO, error)
}
type PaymentService interface {
	CreateByOrderID(ctx context.Context, orderID string) (PaymentDTO, error)
}
type ProductDTO struct {
	ID            string
	Price         float64
	PreparingTime uint
}
type ProductOrderDTO struct{}
type PaymentDTO struct{ QrCode string }

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
	useCase := usecases.Build(orderGateway)

	var productIds []string
	for _, item := range orderDTO.Products {
		productIds = append(productIds, item.ProductID)
	}
	products, err := c.ProductService.FindByIDs(ctx, productIds)
	if err != nil {
		return "", err
	}
	if len(products) != len(orderDTO.Products) {
		return "", err
	}

	totalPrice := 0.0
	totalPreparingTime := uint(0)
	for _, item := range orderDTO.Products {
		for _, product := range products {
			if product.ID == item.ProductID {
				totalPrice += product.Price * float64(item.Quantity)
				totalPreparingTime += product.PreparingTime * uint(item.Quantity)
			}
		}
	}

	order := entity.Order{
		CustomerID:    orderDTO.CustomerID,
		Status:        enum.OrderStatusAwaitingPayment,
		Price:         totalPrice,
		PreparingTime: totalPreparingTime,
	}
	createdOrder, err := useCase.CreateOrder(ctx, order)
	if err != nil {
		return "", err
	}

	productOrders, _ := c.ProductOrderService.BuildBulkFromOrderAndProducts(createdOrder.Entity.ID, orderDTO.Products, products)
	_, err = c.ProductOrderService.CreateBulk(ctx, productOrders)
	if err != nil {
		return "", err
	}

	payment, err := c.PaymentService.CreateByOrderID(ctx, createdOrder.Entity.ID)
	if err != nil {
		return "", err
	}

	return payment.QrCode, nil
}

func (c *Controller) GetAll(ctx context.Context) ([]dto.OrderDAO, error) {
	orderGateway := gateway.Build(c.OrderDatasource)
	useCase := usecases.Build(orderGateway)
	presenter := presenter.Build()

	orders, err := useCase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return presenter.FromEntityListToDAOList(orders), nil
}

func (c *Controller) GetPanel(ctx context.Context, status []string) ([]dto.OrderDAO, error) {
	orderGateway := gateway.Build(c.OrderDatasource)
	useCase := usecases.Build(orderGateway)
	presenter := presenter.Build()

	orders, err := useCase.GetPanel(ctx, status)
	if err != nil {
		return nil, err
	}
	return presenter.FromEntityListToDAOList(orders), nil
}

func (c *Controller) FindByID(ctx context.Context, id string) (dto.OrderDAO, error) {
	orderGateway := gateway.Build(c.OrderDatasource)
	useCase := usecases.Build(orderGateway)
	presenter := presenter.Build()

	order, err := useCase.FindByID(ctx, id)
	if err != nil {
		return dto.OrderDAO{}, err
	}
	return presenter.FromEntityToDAO(order), nil
}

func (c *Controller) Update(ctx context.Context, orderDTO dto.OrderDAO) (dto.OrderDAO, error) {
	orderGateway := gateway.Build(c.OrderDatasource)
	useCase := usecases.Build(orderGateway)
	presenter := presenter.Build()

	order := dto.FromOrderDAO(orderDTO)
	updated, err := useCase.Update(ctx, order)
	if err != nil {
		return dto.OrderDAO{}, err
	}
	return presenter.FromEntityToDAO(updated), nil
}
