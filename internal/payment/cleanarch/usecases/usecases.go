package usecases

import (
	"context"
	"fmt"

	orderenum "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity/enum"
	orderuc "github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/usecases"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/entity/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/cleanarch/gateway"
	productuc "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/usecases"
	productorderuc "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/usecases"
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/cleanarch/entities"
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/cleanarch/external"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type UseCases struct {
	paymentGateway      *gateway.Gateway
	qrCodeProvider      external.QRCodeProvider
	productUseCase      productuc.UseCases
	productOrderUseCase productorderuc.UseCases
	orderUseCase        orderuc.UseCases
}

func Build(
	paymentGateway *gateway.Gateway,
	qrCodeProvider external.QRCodeProvider,
	productUseCase productuc.UseCases,
	productOrderUseCase productorderuc.UseCases,
	orderUseCase orderuc.UseCases,
) *UseCases {
	return &UseCases{
		paymentGateway:      paymentGateway,
		qrCodeProvider:      qrCodeProvider,
		productUseCase:      productUseCase,
		productOrderUseCase: productOrderUseCase,
		orderUseCase:        orderUseCase,
	}
}

func (u *UseCases) CreateByOrderID(ctx context.Context, orderID string) (entity.Payment, error) {
	productOrders, productOrderErr := u.productOrderUseCase.FindByOrderID(ctx, orderID)
	if productOrderErr != nil {
		return entity.Payment{}, productOrderErr
	}

	var productIDs []string
	for _, po := range productOrders {
		productIDs = append(productIDs, po.ProductID)
	}

	products, productsErr := u.productUseCase.FindByIDs(ctx, productIDs)
	if productsErr != nil {
		return entity.Payment{}, productsErr
	}

	var items []entities.Item
	for _, po := range productOrders {
		for _, product := range products {
			if po.ProductID == product.Id {
				items = append(items, entities.Item{
					ID:          product.Id,
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

	qrCode, qrCodeErr := u.qrCodeProvider.GenerateQRCode(ctx, entities.GenerateQRCodeParams{
		OrderID: orderID,
		Items:   items,
	})
	if qrCodeErr != nil {
		return entity.Payment{}, qrCodeErr
	}

	var payment entity.Payment
	createdPayment, createErr := u.paymentGateway.Create(ctx, payment.Build(orderID, qrCode))
	if createErr != nil {
		return entity.Payment{}, createErr
	}

	return createdPayment, nil
}

func (u *UseCases) CheckPayment(ctx context.Context, requestUrl string) (interface{}, error) {
	if requestUrl == "" {
		return nil, &apperror.ValidationError{Msg: "Request URL is required"}
	}

	response, err := u.qrCodeProvider.CheckPayment(ctx, requestUrl)
	if err != nil {
		return nil, fmt.Errorf("error checking payment: %w", err)
	}

	payment, paymentErr := u.paymentGateway.FindByOrderID(ctx, response.ExternalReference)
	if paymentErr != nil {
		return nil, paymentErr
	}
	if response.OrderStatus == "paid" {
		payment.Status = enum.PaymentStatusApproved
		_, updateErr := u.paymentGateway.Update(ctx, payment)
		if updateErr != nil {
			return nil, updateErr
		}

		order, orderErr := u.orderUseCase.FindByID(ctx, response.ExternalReference)
		if orderErr != nil {
			return nil, orderErr
		}

		order.Status = orderenum.OrderStatusReceived
		_, updateOrderErr := u.orderUseCase.Update(ctx, order)
		if updateOrderErr != nil {
			return nil, updateOrderErr
		}
	}

	return response, nil
}
