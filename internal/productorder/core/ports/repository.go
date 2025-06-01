package ports

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/core/model"
)

type ProductOrderRepository interface {
	CreateBulk(ctx context.Context, orders []model.ProductOrder) (int, error)
	FindByOrderID(ctx context.Context, orderID string) ([]model.ProductOrder, error)
}
