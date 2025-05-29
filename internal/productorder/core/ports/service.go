package ports

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/core/model"
)

type ProductOrderService interface {
	CreateBulk(ctx context.Context, productOrders []model.ProductOrder) (int, error)
}
