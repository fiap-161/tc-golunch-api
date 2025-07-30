package services

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/entity/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/usecases"
	paymentinterfaces "github.com/fiap-161/tech-challenge-fiap161/internal/payment/interfaces"
)

type OrderServiceGateway struct {
	orderUseCase *usecases.UseCases
}

func NewOrderServiceGateway(orderUseCase *usecases.UseCases) paymentinterfaces.OrderService {
	return &OrderServiceGateway{
		orderUseCase: orderUseCase,
	}
}

func (a *OrderServiceGateway) FindByID(ctx context.Context, orderID string) (paymentinterfaces.Order, error) {
	order, err := a.orderUseCase.FindByID(ctx, orderID)
	if err != nil {
		return paymentinterfaces.Order{}, err
	}

	return paymentinterfaces.Order{
		ID:     order.Entity.ID,
		Status: string(order.Status),
	}, nil
}

func (a *OrderServiceGateway) Update(ctx context.Context, order paymentinterfaces.Order) (paymentinterfaces.Order, error) {
	currentOrder, err := a.orderUseCase.FindByID(ctx, order.ID)
	if err != nil {
		return paymentinterfaces.Order{}, err
	}

	currentOrder.Status = enum.OrderStatus(order.Status)

	updatedOrder, updateErr := a.orderUseCase.Update(ctx, currentOrder)
	if updateErr != nil {
		return paymentinterfaces.Order{}, updateErr
	}

	return paymentinterfaces.Order{
		ID:     updatedOrder.Entity.ID,
		Status: string(updatedOrder.Status),
	}, nil
}
