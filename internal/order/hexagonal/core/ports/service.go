package ports

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/hexagonal/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/hexagonal/core/model"
)

type OrderService interface {
	Create(ctx context.Context, orderDTO dto.CreateOrderDTO) (string, error)
	GetAll(ctx context.Context) ([]model.Order, error)
	GetPanel(ctx context.Context) ([]model.Order, error)
	Update(ctx context.Context, id, status string) error
}
