package gateway

import (
	"context"

	orderports "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/ports"
	paymentports "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/ports"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/usecases"
)

type ProductOrderServiceAdapter struct {
	productOrderUseCase *usecases.UseCases
}

func NewProductOrderServiceAdapter(productOrderUseCase *usecases.UseCases) (orderports.ProductOrderService, paymentports.ProductOrderService) {
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
