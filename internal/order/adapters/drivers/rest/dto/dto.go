package dto

import (
	"errors"
)

type CreateOrderDTO struct {
	CustomerID string             `json:"customer_id"`
	Products   []OrderProductInfo `json:"products"`
}

type OrderProductInfo struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func (c *CreateOrderDTO) Validate() error {
	if c.CustomerID == "" {
		return errors.New("customer_id is required")
	}
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
