package ports

import (
	"context"

	productentity "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity"
	productorderentity "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/entity"
)

// ProductService define a interface que o domínio Payment precisa para acessar produtos
type ProductService interface {
	FindByIDs(ctx context.Context, productIDs []string) ([]productentity.Product, error)
}

// ProductOrderService define a interface que o domínio Payment precisa para acessar product orders
type ProductOrderService interface {
	FindByOrderID(ctx context.Context, orderID string) ([]productorderentity.ProductOrder, error)
}

// OrderService define a interface que o domínio Payment precisa para acessar orders
type OrderService interface {
	FindByID(ctx context.Context, orderID string) (Order, error)
	Update(ctx context.Context, order Order) (Order, error)
}

// Order define o contrato que o domínio Payment precisa para trabalhar com orders
type Order struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}
