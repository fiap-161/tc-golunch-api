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

func (s *Service) Create(_ context.Context, orderDTO dto.CreateOrderDTO) (string, error) {
	var productIds []string
	for _, item := range orderDTO.Products {
		productIds = append(productIds, item.ProductID)
	}

	//TODO implement productsIds with string
	products, err := s.productRepo.FindByIDs([]uint{})
	if err != nil {
		return "", err
	}
	if len(products) != len(orderDTO.Products) {
		return "", &apperror.NotFoundError{
			Msg: "Some products not found",
		}
	}

	var productOrders = make([]productordermodel.ProductOrder, 0, len(products))
	for _, product := range products {
		for _, item := range orderDTO.Products {
			if string(product.ID) == item.ProductID { //TODO remove this cast
				productOrders = append(productOrders, productordermodel.ProductOrder{
					ProductID: string(product.ID), // TODO remove this cast
					Quantity:  item.Quantity,
					UnitPrice: product.Price,
				})
			}
		}
	}

	_, err = s.productOrderRepo.CreateBulk(context.Background(), productOrders)
	if err != nil {
		return "", err
	}

	var order model.Order
	order = order.FromDTO(orderDTO)

	createdOrder, err := s.orderRepo.Create(nil, order)
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
