package model

import "github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"

type ProductOrder struct {
	entity.Entity
	ProductID string  `json:"product_id"`
	OrderID   string  `json:"order_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}
