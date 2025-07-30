package services

import (
	"context"

	orderinterfaces "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/interfaces"
	paymentinterfaces "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/interfaces"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/usecases"
)

type ProductOrderServiceAdapter struct {
	productOrderUseCase *usecases.UseCases
}

func NewProductOrderServiceAdapter(productOrderUseCase *usecases.UseCases) (
	orderinterfaces.ProductOrderService,
	paymentinterfaces.ProductOrderService,
) {
	adapter := &ProductOrderServiceAdapter{
		productOrderUseCase: productOrderUseCase,
	}
	return adapter, adapter
}

func (a *ProductOrderServiceAdapter) CreateBulk(ctx context.Context, productOrders []entity.ProductOrder) (int, error) {
	return a.productOrderUseCase.CreateBulk(ctx, productOrders)
}

func (a *ProductOrderServiceAdapter) FindByOrderID(ctx context.Context, orderID string) ([]entity.ProductOrder, error) {
	return a.productOrderUseCase.FindByOrderID(ctx, orderID)
}
