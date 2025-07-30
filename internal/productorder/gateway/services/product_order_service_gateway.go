package services

import (
	"context"

	orderinterfaces "github.com/fiap-161/tech-challenge-fiap161/internal/order/interfaces"
	paymentinterfaces "github.com/fiap-161/tech-challenge-fiap161/internal/payment/interfaces"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/usecases"
)

type ProductOrderServiceGateway struct {
	productOrderUseCase *usecases.UseCases
}

func NewProductOrderServiceGateway(productOrderUseCase *usecases.UseCases) (
	orderinterfaces.ProductOrderService,
	paymentinterfaces.ProductOrderService,
) {
	adapter := &ProductOrderServiceGateway{
		productOrderUseCase: productOrderUseCase,
	}
	return adapter, adapter
}

func (a *ProductOrderServiceGateway) CreateBulk(ctx context.Context, productOrders []entity.ProductOrder) (int, error) {
	return a.productOrderUseCase.CreateBulk(ctx, productOrders)
}

func (a *ProductOrderServiceGateway) FindByOrderID(ctx context.Context, orderID string) ([]entity.ProductOrder, error) {
	return a.productOrderUseCase.FindByOrderID(ctx, orderID)
}
