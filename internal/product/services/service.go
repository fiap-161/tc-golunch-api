package services

import (
	"errors"

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

func (s *ProductService) GetAll(category string) ([]model.Product, error) {
	list, err := s.productRepo.GetAll(category)

	if err != nil {
		return nil, &appErrors.InternalError{Msg: "Error querying table"}
	}

	return list, nil
}

func (s *ProductService) Update(product model.Product, id uint) (model.Product, error) {
	product, err := s.productRepo.Update(id, product)

	if err != nil {
		var notFoundErr *appErrors.NotFoundError
		if errors.As(err, &notFoundErr) {
			return model.Product{}, notFoundErr
		}
		return model.Product{}, &appErrors.InternalError{Msg: "Unexpected error"}
	}

	return product, nil
}

func (s *ProductService) FindByID(id uint) (model.Product, error) {
	product, err := s.productRepo.FindByID(id)

	if err != nil {
		var notFoundErr *appErrors.NotFoundError
		if errors.As(err, &notFoundErr) {
			return model.Product{}, notFoundErr
		}
		return model.Product{}, &appErrors.InternalError{Msg: "Unexpected error"}
	}

	return product, nil
}

func (s *ProductService) Delete(id uint) error {
	err := s.productRepo.Delete(id)

	if err != nil {
		var notFoundErr *appErrors.NotFoundError
		if errors.As(err, &notFoundErr) {
			return notFoundErr
		}
		return &appErrors.InternalError{Msg: "Unexpected error"}
	}

	return nil
}
