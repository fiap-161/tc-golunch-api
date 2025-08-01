package interfaces

import (
	"context"

	paymententity "github.com/fiap-161/tech-challenge-fiap161/internal/payment/entity"
	productentity "github.com/fiap-161/tech-challenge-fiap161/internal/product/entity"
	productorderentity "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/entity"
)

type ProductService interface {
	FindByIDs(ctx context.Context, productIDs []string) ([]productentity.Product, error)
}

type ProductOrderService interface {
	CreateBulk(ctx context.Context, productOrders []productorderentity.ProductOrder) (int, error)
}

type PaymentService interface {
	CreateByOrderID(ctx context.Context, orderID string) (*paymententity.Payment, error)
}
