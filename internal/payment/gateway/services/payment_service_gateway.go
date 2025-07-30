package services

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/interfaces"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/usecases"
)

type PaymentServiceGateway struct {
	paymentUseCase *usecases.UseCases
}

func NewPaymentServiceGateway(paymentUseCase *usecases.UseCases) interfaces.PaymentService {
	return &PaymentServiceGateway{
		paymentUseCase: paymentUseCase,
	}
}

func (a *PaymentServiceGateway) CreateByOrderID(ctx context.Context, orderID string) (*entity.Payment, error) {
	payment, err := a.paymentUseCase.CreateByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
