package usecases

import (
	"context"
	"errors"
	"fmt"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/gateway"
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

type UseCases struct {
	OrderGateway        *gateway.Gateway
	ProductService      ProductService
	ProductOrderService ProductOrderService
	PaymentService      PaymentService
}

func Build(orderGateway *gateway.Gateway, productService ProductService, productOrderService ProductOrderService, paymentService PaymentService) *UseCases {
	return &UseCases{
		OrderGateway:        orderGateway,
		ProductService:      productService,
		ProductOrderService: productOrderService,
		PaymentService:      paymentService,
	}
}

func (u *UseCases) CreateCompleteOrder(ctx context.Context, orderDTO dto.CreateOrderDTO) (string, error) {
	products, err := u.validateProducts(ctx, orderDTO.Products)
	if err != nil {
		return "", err
	}

	totalPrice, totalPreparingTime := u.calculateOrderTotals(orderDTO.Products, products)

	order := entity.Order{
		CustomerID:    orderDTO.CustomerID,
		Status:        enum.OrderStatusAwaitingPayment,
		Price:         totalPrice,
		PreparingTime: totalPreparingTime,
	}

	createdOrder, err := u.CreateOrder(ctx, order)
	if err != nil {
		return "", fmt.Errorf("error when creating order: %w", err)
	}

	err = u.createProductOrderRelations(ctx, createdOrder.ID, orderDTO.Products, products)
	if err != nil {
		return "", fmt.Errorf("error when creating product-order relations: %w", err)
	}

	payment, err := u.PaymentService.CreateByOrderID(ctx, createdOrder.ID)
	if err != nil {
		return "", fmt.Errorf("error when creating payment: %w", err)
	}

	return payment.QrCode, nil
}

func (u *UseCases) validateProducts(ctx context.Context, orderProducts []dto.OrderProductInfo) ([]dto.ProductDTO, error) {
	if len(orderProducts) == 0 {
		return nil, errors.New("order need to have at least one product")
	}

	var productIds []string
	for _, item := range orderProducts {
		if item.Quantity <= 0 {
			return nil, errors.New("quantity need to be greater than zero")
		}
		productIds = append(productIds, item.ProductID)
	}

	products, err := u.ProductService.FindByIDs(ctx, productIds)
	if err != nil {
		return nil, fmt.Errorf("error when search products: %w", err)
	}

	if len(products) != len(orderProducts) {
		return nil, errors.New("products need to have at least one product")
	}

	return products, nil
}

func (u *UseCases) calculateOrderTotals(orderProducts []dto.OrderProductInfo, products []dto.ProductDTO) (float64, uint) {
	totalPrice := 0.0
	totalPreparingTime := uint(0)

	for _, item := range orderProducts {
		for _, product := range products {
			if product.ID == item.ProductID {
				itemPrice := product.Price * float64(item.Quantity)
				itemPreparingTime := product.PreparingTime * uint(item.Quantity)

				totalPrice += itemPrice
				totalPreparingTime += itemPreparingTime
				break
			}
		}
	}

	return totalPrice, totalPreparingTime
}

func (u *UseCases) createProductOrderRelations(ctx context.Context, orderID string, orderProducts []dto.OrderProductInfo, products []dto.ProductDTO) error {
	productOrders, err := u.ProductOrderService.BuildBulkFromOrderAndProducts(orderID, orderProducts, products)
	if err != nil {
		return fmt.Errorf("error when building product-order relationships: %w", err)
	}

	_, err = u.ProductOrderService.CreateBulk(ctx, productOrders)
	if err != nil {
		return fmt.Errorf("error persisting product-order relationships: %w", err)
	}

	return nil
}

func (u *UseCases) CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error) {
	return u.OrderGateway.Create(ctx, order)
}

func (u *UseCases) GetAll(ctx context.Context) ([]entity.Order, error) {
	return u.OrderGateway.GetAll(ctx)
}

func (u *UseCases) GetPanel(ctx context.Context, status []string) ([]entity.Order, error) {
	return u.OrderGateway.GetPanel(ctx, status)
}

func (u *UseCases) FindByID(ctx context.Context, id string) (entity.Order, error) {
	return u.OrderGateway.FindByID(ctx, id)
}

func (u *UseCases) Update(ctx context.Context, order entity.Order) (entity.Order, error) {
	return u.OrderGateway.Update(ctx, order)
}
