package controllers

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/entity/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/external/datasource"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/gateway"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/presenter"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/usecases"
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/cleanarch/external"
)

type Controller struct {
	PaymentDatasource   datasource.DataSource
	QRCodeProvider      external.QRCodeProvider
	ProductService      ProductUsecase
	ProductOrderService ProductOrderService
	OrderService        OrderService
}

func Build(paymentDatasource datasource.DataSource, qrCodeProvider QRCodeProvider, productService ProductService, productOrderService ProductOrderService, orderService OrderService) *Controller {
	return &Controller{
		PaymentDatasource:   paymentDatasource,
		QRCodeProvider:      qrCodeProvider,
		ProductService:      productService,
		ProductOrderService: productOrderService,
		OrderService:        orderService,
	}
}

func (c *Controller) CreateByOrderID(ctx context.Context, orderID string) (dto.PaymentResponseDTO, error) {
	paymentGateway := gateway.Build(c.PaymentDatasource)
	useCase := usecases.Build(paymentGateway, c.QRCodeProvider, c.ProductService, c.ProductOrderService, c.OrderService)
	presenter := presenter.Build()

	payment, err := useCase.CreateByOrderID(ctx, orderID)
	if err != nil {
		return dto.PaymentResponseDTO{}, err
	}

	return presenter.FromEntityToResponseDTO(payment), nil
}

func (c *Controller) CheckPayment(ctx context.Context, requestUrl string) (interface{}, error) {
	paymentGateway := gateway.Build(c.PaymentDatasource)
	useCase := usecases.Build(paymentGateway, c.QRCodeProvider, c.ProductService, c.ProductOrderService, c.OrderService)

	return useCase.CheckPayment(ctx, requestUrl)
}
