package ports

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/core/model"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment model.Payment) (model.Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (model.Payment, error)
	Update(ctx context.Context, payment model.Payment) (model.Payment, error)
}
