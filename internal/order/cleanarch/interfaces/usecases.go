package interfaces

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity"
)

type UseCase interface {
	CreateCompleteOrder(ctx context.Context, orderDTO dto.CreateOrderDTO) (string, error)
	CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error)
	GetAll(ctx context.Context) ([]entity.Order, error)
	GetPanel(ctx context.Context, status []string) ([]entity.Order, error)
	FindByID(ctx context.Context, orderID string) (entity.Order, error)
	Update(ctx context.Context, order entity.Order) (entity.Order, error)
}
