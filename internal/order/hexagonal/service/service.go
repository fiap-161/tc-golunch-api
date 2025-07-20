package service

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/hexagonal/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/hexagonal/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/hexagonal/core/ports"

	paymentport "github.com/fiap-161/tech-challenge-fiap161/internal/payment/core/ports"
	productcontroller "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/controller"
	productordercontroller "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/controller"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type Service struct {
	orderRepo              ports.OrderRepository
	productController      productcontroller.Controller
	productOrderController productordercontroller.Controller
	paymentService         paymentport.PaymentService
}

func New(
	orderRepo ports.OrderRepository,
	productController productcontroller.Controller,
	productOrderController productordercontroller.Controller,
	paymentService paymentport.PaymentService,
) ports.OrderService {
	return &Service{
		orderRepo:              orderRepo,
		productController:      productController,
		productOrderController: productOrderController,
		paymentService:         paymentService,
	}
}

func (s *Service) Create(ctx context.Context, orderDTO dto.CreateOrderDTO) (string, error) {
	var productIds []string
	for _, item := range orderDTO.Products {
		productIds = append(productIds, item.ProductID)
	}

	products, findErr := s.productController.FindByIDs(ctx, productIds)
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

	productOrders, _ := s.productOrderController.BuildBulkFromOrderAndProducts(createdOrder.ID, orderDTO.Products, products)
	_, createBulkErr := s.productOrderController.CreateBulk(ctx, productOrders)
	if createBulkErr != nil {
		return "", createBulkErr
	}

	payment, paymentErr := s.paymentService.CreateByOrderID(ctx, createdOrder.ID)
	if paymentErr != nil {
		return "", paymentErr
	}

	return payment.QrCode, nil
}

func (s *Service) Update(ctx context.Context, id, status string) error {
	order, err := s.orderRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status == model.OrderStatusAwaitingPayment {
		return &apperror.ValidationError{Msg: "order status must be different from awaiting payment"}
	}

	_, err = s.orderRepo.Update(ctx, order.BuildUpdate(model.OrderStatus(status)))
	return err
}

func (s *Service) GetAll(ctx context.Context) ([]model.Order, error) {
	order, err := s.orderRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) GetPanel(ctx context.Context) ([]model.Order, error) {
	orders, err := s.orderRepo.GetPanel(ctx, model.OrderPanelStatus)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
