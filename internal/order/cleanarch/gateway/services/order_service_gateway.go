package services

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/usecases"
	paymentports "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/ports"
)

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
	currentOrder, err := a.orderUseCase.FindByID(ctx, order.ID)
	if err != nil {
		return paymentports.Order{}, err
	}

	currentOrder.Status = enum.OrderStatus(order.Status)

	updatedOrder, updateErr := a.orderUseCase.Update(ctx, currentOrder)
	if updateErr != nil {
		return paymentports.Order{}, updateErr
	}

	return paymentports.Order{
		ID:     updatedOrder.Entity.ID,
		Status: string(updatedOrder.Status),
	}, nil
}
