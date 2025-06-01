package model

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
	"github.com/google/uuid"
)

type ProductOrder struct {
	entity.Entity
	ProductID string  `json:"product_id"`
	OrderID   string  `json:"order_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

func BuildBulkFromOrderAndProducts(orderID string, orderProductInfo []dto.OrderProductInfo, products []model.Product) []ProductOrder {

	var productOrders = make([]ProductOrder, 0, len(products))
	for _, product := range products {
		for _, item := range orderProductInfo {
			if string(product.ID) == item.ProductID { //TODO remove this cast
				productOrders = append(productOrders, ProductOrder{
					Entity: entity.Entity{
						ID: uuid.NewString(),
					},
					OrderID:   orderID,
					ProductID: string(product.ID), // TODO remove this cast
					Quantity:  item.Quantity,
					UnitPrice: product.Price,
				})
			}
		}
	}

	return productOrders
}
