package ports

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/core/model/enum"
)

type ProductService interface {
	Create(context.Context, model.Product) (model.Product, error)
	ListCategories(context.Context) []enum.Category
	GetAll(context.Context, enum.Category) ([]model.Product, error)
	Update(context.Context, model.Product, string) (model.Product, error)
	FindByID(context.Context, string) (model.Product, error)
	Delete(context.Context, string) error
}
