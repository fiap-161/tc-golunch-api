package ports

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
)

type ProductRepository interface {
	Create(context.Context, model.Product) (model.Product, error)
	GetAll(context.Context, uint) ([]model.Product, error)
	Update(context.Context, string, model.Product) (model.Product, error)
	FindByID(context.Context, string) (model.Product, error)
	FindByIDs(context.Context, []string) ([]model.Product, error)
	Delete(context.Context, string) error
}
