package dto

type CreateOrderDTO struct {
	CustomerID string `json:"customer_id"`
	Products   []uint `json:"products"`
}
