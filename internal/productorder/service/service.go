package service

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/core/model"
	productorderport "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/core/ports"
)

type Service struct {
	productOrderRepo productorderport.ProductOrderRepository
}

func New(productOrderRepo productorderport.ProductOrderRepository) productorderport.ProductOrderService {
	return &Service{
		productOrderRepo: productOrderRepo,
	}
}

func (s *Service) CreateBulk(ctx context.Context, productOrders []model.ProductOrder) (int, error) {
	if len(productOrders) == 0 {
		return 0, nil
	}

	createdCount, err := s.productOrderRepo.CreateBulk(ctx, productOrders)
	if err != nil {
		return 0, err
	}

	return createdCount, nil
}
