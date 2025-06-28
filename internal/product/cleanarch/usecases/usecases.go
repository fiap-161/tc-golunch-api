package usecases

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/gateway"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type UseCases struct {
	ProductGateway gateway.Gateway
}

func Build(productGateway gateway.Gateway) *UseCases {
	return &UseCases{ProductGateway: productGateway}
}

func (u *UseCases) CreateProduct(ctx context.Context, productDTO dto.ProductRequestDTO) (entity.Product, error) {
	var product entity.Product
	product = product.FromRequestDTO(productDTO)
	isValidCategory := enum.IsValidCategory(string(product.Category))

	if !isValidCategory {
		return entity.Product{}, &apperror.ValidationError{Msg: "Invalid category"}
	}

	saved, err := u.ProductGateway.Create(ctx, product.Build())
	if err != nil {
		return entity.Product{}, &apperror.InternalError{Msg: err.Error()}
	}

	return saved, nil
}

func (u *UseCases) ListCategories(ctx context.Context) []enum.Category {
	return enum.GetAllCategories()
}
