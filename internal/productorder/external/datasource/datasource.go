package datasource

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/dto"
)

type DataSource interface {
	CreateBulk(ctx context.Context, orders []dto.ProductOrderDAO) (int, error)
	FindByOrderID(ctx context.Context, orderID string) ([]dto.ProductOrderDAO, error)
}
