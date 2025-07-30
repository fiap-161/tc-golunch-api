package services

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/interfaces"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/usecases"
)

type ProductServiceGateway struct {
	productUseCase *usecases.UseCases
}

func NewProductServiceGateway(productUseCase *usecases.UseCases) interfaces.ProductService {
	return &ProductServiceGateway{
		productUseCase: productUseCase,
	}
}

func (a *ProductServiceGateway) FindByIDs(ctx context.Context, productIDs []string) ([]entity.Product, error) {
	return a.productUseCase.FindByIDs(ctx, productIDs)
}
