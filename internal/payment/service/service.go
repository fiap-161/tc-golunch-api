package service

import (
	"context"
	"fmt"
	ordermodel "github.com/fiap-161/tech-challenge-fiap161/internal/order/core/model"

	orderport "github.com/fiap-161/tech-challenge-fiap161/internal/order/core/ports"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/core/ports"
	productport "github.com/fiap-161/tech-challenge-fiap161/internal/product/core/ports"
	productorderport "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/core/ports"
	qrcodedto "github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/core/dto"
	qrcodeports "github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/core/ports"
)

type service struct {
	qrCodeClient     qrcodeports.QRCodeProvider
	orderRepo        orderport.OrderRepository
	paymentRepo      ports.PaymentRepository
	productOrderRepo productorderport.ProductOrderRepository
	productRepo      productport.ProductRepository
}

func New(
	qrCodeClient qrcodeports.QRCodeProvider,
	orderRepo orderport.OrderRepository,
	paymentRepo ports.PaymentRepository,
	productOrderRepo productorderport.ProductOrderRepository,
	productRepo productport.ProductRepository,
) ports.PaymentService {
	return &service{
		qrCodeClient:     qrCodeClient,
		orderRepo:        orderRepo,
		paymentRepo:      paymentRepo,
		productOrderRepo: productOrderRepo,
		productRepo:      productRepo,
	}
}

func (s *service) CreateByOrderID(ctx context.Context, orderID string) (model.Payment, error) {
	productOrders, productOrderErr := s.productOrderRepo.FindByOrderID(ctx, orderID)
	if productOrderErr != nil {
		return model.Payment{}, productOrderErr
	}

	var productIDs []string
	for _, po := range productOrders {
		productIDs = append(productIDs, po.ProductID)
	}

	products, productErr := s.productRepo.FindByIDs(ctx, productIDs)
	if productErr != nil {
		return model.Payment{}, productErr
	}

	var items []qrcodedto.Item
	for _, po := range productOrders {
		for _, product := range products {
			if po.ProductID == product.ID {
				items = append(items, qrcodedto.Item{
					ID:          product.ID,
					Name:        product.Name,
					Price:       product.Price,
					Description: product.Description,
					Quantity:    po.Quantity,
					Amount:      product.Price * float64(po.Quantity),
				})
				break
			}
		}
	}

	qrCode, qrCodeErr := s.qrCodeClient.GenerateQRCode(ctx, qrcodedto.GenerateQRCodeParams{
		OrderID: orderID,
		Items:   items,
	})
	if qrCodeErr != nil {
		return model.Payment{}, qrCodeErr
	}

	var payment model.Payment
	createdPayment, createErr := s.paymentRepo.Create(ctx, payment.Build(orderID, qrCode))
	if createErr != nil {
		return model.Payment{}, createErr
	}

	return createdPayment, nil
}

func (s *service) CheckPayment(ctx context.Context, requestUrl string) (any, error) {
	response, err := s.qrCodeClient.CheckPayment(ctx, requestUrl)
	if err != nil {
		return nil, err
	}

	fmt.Println(response)

	payment, paymentErr := s.paymentRepo.FindByOrderID(ctx, response.ExternalReference)
	if paymentErr != nil {
		return nil, paymentErr
	}
	if response.OrderStatus == "paid" {
		payment.Status = model.Approved
		_, updateErr := s.paymentRepo.Update(ctx, payment)
		if updateErr != nil {
			return nil, updateErr
		}

		order, orderErr := s.orderRepo.FindByID(ctx, response.ExternalReference)
		if orderErr != nil {
			return nil, orderErr
		}

		order.Status = ordermodel.OrderStatusReceived
		_, updateOrderErr := s.orderRepo.Update(ctx, order)
		if updateOrderErr != nil {
			return nil, updateOrderErr
		}
	}

	// TODO verify possible statuses
	//if verifyResponse.OrderStatus == "CANCELED" {
	//	payment.Status = model.Rejected
	//}

	return response, nil
}
