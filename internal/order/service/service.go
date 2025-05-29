package service

import (
	"context"
	"encoding/json"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/core/model"
	orderport "github.com/fiap-161/tech-challenge-fiap161/internal/order/core/ports"
	productport "github.com/fiap-161/tech-challenge-fiap161/internal/product/core/ports"
)

type Service struct {
	orderRepo   orderport.OrderRepository
	productRepo productport.ProductRepository
}

func New(orderRepo orderport.OrderRepository, productRepo productport.ProductRepository) orderport.OrderServicePort {
	return &Service{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *Service) Create(ctx context.Context, orderDTO dto.CreateOrderDTO) (string, error) {
	var order model.Order
	order = order.FromDTO(orderDTO)

	productsIds := orderDTO.Products
	products, err := s.productRepo.FindByIDs(productsIds)
	if err != nil {
		return "", err
	}

	var totalPrice float64
	var estimatedTime uint
	for _, product := range products {
		totalPrice += product.Price
		estimatedTime += product.PreparingTime
	}

	productsJSON, err := json.Marshal(orderDTO.Products)
	if err != nil {
		return "", err
	}

	createdOrder, err := s.orderRepo.Create(ctx, order.Build(totalPrice, estimatedTime, productsJSON))
	if err != nil {
		return "", err
	}

	return createdOrder.ID, nil
}

func (s *Service) GetAll(ctx context.Context) ([]model.Order, error) {
	order, err := s.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return order, nil
}
