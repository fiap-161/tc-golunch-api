package datasource

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity"
)

// TODO: Como é mundo externo, precisa ser uma DTO
type DataSource interface {
	Create(context.Context, entity.Product) (entity.Product, error)
	GetAll(context.Context, uint) ([]entity.Product, error)
	Update(context.Context, string, entity.Product) (entity.Product, error)
	FindByID(context.Context, string) (entity.Product, error)
	FindByIDs(context.Context, []string) ([]entity.Product, error)
	Delete(context.Context, string) error
}
