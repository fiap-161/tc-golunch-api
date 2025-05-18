package ports

import "github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"

type ProductRepository interface {
	Create(product model.Product) (model.Product, error)
	GetAll(category string) ([]model.Product, error)
	Update(id uint, updated model.Product) (model.Product, error)
	FindById(id uint) (model.Product, error)
	Delete(id uint) error
}
