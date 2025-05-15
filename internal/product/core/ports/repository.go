package ports

import "github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"

type ProductRepository interface {
	Create(product model.Product) (model.Product, error)
}
