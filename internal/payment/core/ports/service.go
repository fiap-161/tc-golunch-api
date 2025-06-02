package ports

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/core/model"
)

type PaymentService interface {
	CreateByOrderID(ctx context.Context, orderID string) (model.Payment, error)
	CheckPayment(ctx context.Context, requestUrl string) (any, error)
}
