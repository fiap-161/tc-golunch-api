package services

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/interfaces"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/usecases"
)

type ProductServiceAdapter struct {
	productUseCase *usecases.UseCases
}

func NewProductServiceAdapter(productUseCase *usecases.UseCases) interfaces.ProductService {
	return &ProductServiceAdapter{
		productUseCase: productUseCase,
	}
}

func (a *ProductServiceAdapter) FindByIDs(ctx context.Context, productIDs []string) ([]entity.Product, error) {
	return a.productUseCase.FindByIDs(ctx, productIDs)
}
