package ports

import (
	"context"

	paymententity "github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/entity"
	productentity "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity"
	productorderentity "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/entity"
)

// ProductService define a interface que o domínio Order precisa para acessar produtos
type ProductService interface {
	FindByIDs(ctx context.Context, productIDs []string) ([]productentity.Product, error)
}

// ProductOrderService define a interface que o domínio Order precisa para criar product orders
type ProductOrderService interface {
	CreateBulk(ctx context.Context, productOrders []productorderentity.ProductOrder) (int, error)
}

// PaymentService define a interface que o domínio Order precisa para criar pagamentos
type PaymentService interface {
	CreateByOrderID(ctx context.Context, orderID string) (*paymententity.Payment, error)
}
