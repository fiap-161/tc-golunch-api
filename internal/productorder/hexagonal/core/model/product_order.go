package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
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
			if product.ID == item.ProductID {
				productOrders = append(productOrders, ProductOrder{
					Entity: entity.Entity{
						ID:        uuid.NewString(),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					OrderID:   orderID,
					ProductID: product.ID,
					Quantity:  item.Quantity,
					UnitPrice: product.Price,
				})
			}
		}
	}

	return productOrders
}
