package ports

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"
)

type ProductRepository interface {
	Create(context.Context, model.Product) (model.Product, error)
	GetAll(context.Context, enum.Category) ([]model.Product, error)
	Update(context.Context, string, model.Product) (model.Product, error)
	FindByID(context.Context, string) (model.Product, error)
	FindByIDs(context.Context, []string) ([]model.Product, error)
	Delete(context.Context, string) error
}
