package service

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/core/model"
	orderport "github.com/fiap-161/tech-challenge-fiap161/internal/order/core/ports"
	productport "github.com/fiap-161/tech-challenge-fiap161/internal/product/core/ports"
	productordermodel "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/core/model"
	productorderport "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/core/ports"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type Service struct {
	orderRepo        orderport.OrderRepository
	productRepo      productport.ProductRepository
	productOrderRepo productorderport.ProductOrderRepository
}

func New(orderRepo orderport.OrderRepository, productRepo productport.ProductRepository) orderport.OrderService {
	return &Service{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *Service) Create(ctx context.Context, orderDTO dto.CreateOrderDTO) (string, error) {
	var productIds []string
	for _, item := range orderDTO.Products {
		productIds = append(productIds, item.ProductID)
	}

	//TODO implement productsIds with string
	products, findErr := s.productRepo.FindByIDs([]uint{})
	if findErr != nil {
		return "", findErr
	}
	if len(products) != len(orderDTO.Products) {
		return "", &apperror.NotFoundError{
			Msg: "some products not found",
		}
	}

	var order model.Order
	order = order.FromDTO(orderDTO, products)
	createdOrder, createErr := s.orderRepo.Create(ctx, order.Build())
	if createErr != nil {
		return "", createErr
	}

	productOrders := productordermodel.BuildBulkFromOrderAndProducts(createdOrder.ID, orderDTO.Products, products)

	_, createBulkErr := s.productOrderRepo.CreateBulk(ctx, productOrders)
	if createBulkErr != nil {
		return "", createBulkErr
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
