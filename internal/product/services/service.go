package services

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/ports"
	appErrors "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type ProductService struct {
	productRepo ports.ProductRepository
}

func NewProductService(productRepository ports.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepository}
}

func (s *ProductService) Create(product model.Product) (model.Product, error) {
	isValidCategory := enum.IsValidCategory(uint(product.Category))

	if !isValidCategory {
		return model.Product{}, &appErrors.ValidationError{Msg: "Invalid category"}
	}

	savedProduct, err := s.productRepo.Create(product)
	if err != nil {
		return model.Product{}, &appErrors.InternalError{Msg: "Error saving product"}
	}

	return savedProduct, nil
}

func (s *ProductService) ListCategories() []enum.CategoryDTO {
	return enum.GetAllCategories()
}
