package dto

import "github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"

type ProductOrderDAO struct {
	entity.Entity
	ProductID string  `json:"product_id"`
	OrderID   string  `json:"order_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type ProductOrderRequestDTO struct {
	ProductID string  `json:"product_id"`
	OrderID   string  `json:"order_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type ProductOrderResponseDTO struct {
	ID        string  `json:"id"`
	ProductID string  `json:"product_id"`
	OrderID   string  `json:"order_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}
