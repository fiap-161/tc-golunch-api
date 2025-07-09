package service

import (
	"context"
	"errors"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/core/model/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/core/ports"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type Service struct {
	repo ports.ProductRepository
}

func New(repo ports.ProductRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, product model.Product) (model.Product, error) {
	isValidCategory := enum.IsValidCategory(string(product.Category))

	if !isValidCategory {
		return model.Product{}, &apperror.ValidationError{Msg: "Invalid category"}
	}

	saved, err := s.repo.Create(ctx, product.Build())
	if err != nil {
		return model.Product{}, &apperror.InternalError{Msg: err.Error()}
	}

	return saved, nil
}

func (s *Service) ListCategories(_ context.Context) []enum.Category {
	return enum.GetAllCategories()
}

func (s *Service) GetAll(ctx context.Context, category enum.Category) ([]model.Product, error) {
	list, err := s.repo.GetAll(ctx, category)

	if err != nil {
		return nil, &apperror.InternalError{Msg: err.Error()}
	}

	return list, nil
}

func (s *Service) Update(ctx context.Context, product model.Product, id string) (model.Product, error) {
	updated, err := s.repo.Update(ctx, id, product)

	if err != nil {
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			return model.Product{}, notFoundErr
		}
		return model.Product{}, &apperror.InternalError{Msg: err.Error()}
	}

	return updated, nil
}

func (s *Service) FindByID(ctx context.Context, id string) (model.Product, error) {
	product, err := s.repo.FindByID(ctx, id)

	if err != nil {
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			return model.Product{}, notFoundErr
		}
		return model.Product{}, &apperror.InternalError{Msg: "Unexpected error"}
	}

	return product, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.repo.Delete(ctx, id)

	if err != nil {
		var notFoundErr *apperror.NotFoundError
		if errors.As(err, &notFoundErr) {
			return notFoundErr
		}
		return &apperror.InternalError{Msg: err.Error()}
	}

	return nil
}
