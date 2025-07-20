package usecases

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/gateway"
)

type UseCases struct {
	OrderGateway *gateway.Gateway
}

func Build(orderGateway *gateway.Gateway) *UseCases {
	return &UseCases{OrderGateway: orderGateway}
}

func (u *UseCases) CreateOrder(ctx context.Context, order entity.Order) (entity.Order, error) {
	return u.OrderGateway.Create(ctx, order)
}

func (u *UseCases) GetAll(ctx context.Context) ([]entity.Order, error) {
	return u.OrderGateway.GetAll(ctx)
}

func (u *UseCases) GetPanel(ctx context.Context, status []string) ([]entity.Order, error) {
	return u.OrderGateway.GetPanel(ctx, status)
}

func (u *UseCases) FindByID(ctx context.Context, id string) (entity.Order, error) {
	return u.OrderGateway.FindByID(ctx, id)
}

func (u *UseCases) Update(ctx context.Context, order entity.Order) (entity.Order, error) {
	return u.OrderGateway.Update(ctx, order)
} 