package adapters

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/usecases"
	paymentports "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/ports"
)

// OrderServiceAdapter implementa a interface OrderService do domínio Payment
type OrderServiceAdapter struct {
	orderUseCase *usecases.UseCases
}

func NewOrderServiceAdapter(orderUseCase *usecases.UseCases) paymentports.OrderService {
	return &OrderServiceAdapter{
		orderUseCase: orderUseCase,
	}
}

func (a *OrderServiceAdapter) FindByID(ctx context.Context, orderID string) (paymentports.Order, error) {
	order, err := a.orderUseCase.FindByID(ctx, orderID)
	if err != nil {
		return paymentports.Order{}, err
	}

	return paymentports.Order{
		ID:     order.Entity.ID,
		Status: string(order.Status),
	}, nil
}

func (a *OrderServiceAdapter) Update(ctx context.Context, order paymentports.Order) (paymentports.Order, error) {
	// Por enquanto, retornar a mesma order para evitar o ciclo de import
	// TODO: Implementar atualização adequada quando resolver o ciclo de import
	return order, nil
}
