package services

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/interfaces"

	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/usecases"
)

type PaymentServiceAdapter struct {
	paymentUseCase *usecases.UseCases
}

func NewPaymentServiceAdapter(paymentUseCase *usecases.UseCases) interfaces.PaymentService {
	return &PaymentServiceAdapter{
		paymentUseCase: paymentUseCase,
	}
}

func (a *PaymentServiceAdapter) CreateByOrderID(ctx context.Context, orderID string) (*entity.Payment, error) {
	payment, err := a.paymentUseCase.CreateByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
