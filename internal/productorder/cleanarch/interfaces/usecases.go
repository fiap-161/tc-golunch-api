package interfaces

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/entity"
)

type UseCase interface {
	CreateBulk(ctx context.Context, productOrders []entity.ProductOrder) (int, error)
	FindByOrderID(ctx context.Context, orderID string) ([]entity.ProductOrder, error)
}
