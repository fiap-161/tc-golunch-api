package usecases

import (
	"context"
	"fmt"

	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/gateway"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type UseCases struct {
	ProductOrderGateway gateway.Gateway
}

func Build(productGateway gateway.Gateway) *UseCases {
	return &UseCases{ProductOrderGateway: productGateway}
}

func (u *UseCases) CreateBulk(ctx context.Context, productOrders []entity.ProductOrder) (int, error) {

	for i, po := range productOrders {
		if po.ProductID == "" {
			return 0, &apperror.ValidationError{Msg: fmt.Sprintf("productOrder[%d]: ProductID não pode ser vazio", i)}
		}
		if po.OrderID == "" {
			return 0, &apperror.ValidationError{Msg: fmt.Sprintf("productOrder[%d]: OrderID não pode ser vazio", i)}
		}
		if po.Quantity <= 0 {
			return 0, &apperror.ValidationError{Msg: fmt.Sprintf("productOrder[%d]: Quantity deve ser maior que zero", i)}
		}
		if po.UnitPrice < 0 {
			return 0, &apperror.ValidationError{Msg: fmt.Sprintf("productOrder[%d]: UnitPrice não pode ser negativo", i)}
		}
	}

	length, err := u.ProductOrderGateway.CreateBulk(ctx, productOrders)
	if err != nil {
		return 0, err
	}

	return length, nil
}

func (u *UseCases) FindByOrderID(ctx context.Context, orderId string) ([]entity.ProductOrder, error) {

	productOrderFound, err := u.ProductOrderGateway.FindByOrderID(ctx, orderId)
	if err != nil {
		return []entity.ProductOrder{}, err
	}

	return productOrderFound, nil
}
