package interfaces

import (
	"context"

	productentity "github.com/fiap-161/tech-challenge-fiap161/internal/product/entity"
	productorderentity "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/entity"
)

type ProductService interface {
	FindByIDs(ctx context.Context, productIDs []string) ([]productentity.Product, error)
}

type ProductOrderService interface {
	FindByOrderID(ctx context.Context, orderID string) ([]productorderentity.ProductOrder, error)
}

type OrderService interface {
	FindByID(ctx context.Context, orderID string) (Order, error)
	Update(ctx context.Context, order Order) (Order, error)
}

type Order struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
