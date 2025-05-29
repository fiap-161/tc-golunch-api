package ports

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/core/model"
)

type OrderService interface {
	Create(ctx context.Context, orderDTO dto.CreateOrderDTO) (string, error)
	GetAll(ctx context.Context) ([]model.Order, error)
}
