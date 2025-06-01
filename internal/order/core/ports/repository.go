package ports

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/core/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) (model.Order, error)
	GetAll(ctx context.Context) ([]model.Order, error)
	FindByID(ctx context.Context, id string) (model.Order, error)
	Update(ctx context.Context, order model.Order) (model.Order, error)
}
