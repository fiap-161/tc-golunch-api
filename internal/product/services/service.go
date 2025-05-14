package services

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"
	appErrors "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type ProductService struct {
}

func NewProductService() *ProductService {
	return &ProductService{}
}

func (s *ProductService) Create(product model.Product) (model.Product, error) {
	isValidCategory := enum.IsValidCategory(uint(product.Category))

	if !isValidCategory {
		return model.Product{}, &appErrors.ValidationError{Msg: "Invalid category"}
	}
	return model.Product{
		ID:            1,
		Name:          product.Name,
		Price:         product.Price,
		Description:   product.Description,
		PreparingTime: product.PreparingTime,
		Category:      product.Category,
		ImageURL:      product.ImageURL,
	}, nil
}

func (s *ProductService) ListCategories() []enum.CategoryDTO {
	return enum.GetAllCategories()
}
