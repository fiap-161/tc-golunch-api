package dto

import (
	"errors"
)

type CreateOrderDTO struct {
	CustomerID string             `json:"customer_id"`
	Products   []OrderProductInfo `json:"products"`
}

type UpdateOrderDTO struct {
	Status string `json:"status" binding:"required"`
}

type OrderProductInfo struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderPanelDTO struct {
	Orders []OrderPanelItemDTO `json:"orders"`
}

type OrderPanelItemDTO struct {
	OrderNumber   string `json:"order_number"`
	Status        string `json:"status"`
	PreparingTime uint   `json:"preparing_time"`
}

func (c *CreateOrderDTO) Validate() error {
	if len(c.Products) == 0 {
		return errors.New("at least one product is required")
	}
	for _, v := range c.Products {
		if v.ProductID == "" {
			return errors.New("products must not contain empty values")
		}

		if v.Quantity <= 0 {
			return errors.New("product quantity must be greater than zero")
		}
	}
	return nil
}
