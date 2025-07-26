package ports

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/model"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer model.Customer) (model.Customer, error)
	FindByCPF(ctx context.Context, cpf string) (model.Customer, error)
}
