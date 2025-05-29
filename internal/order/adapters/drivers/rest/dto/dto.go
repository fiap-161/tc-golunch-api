package dto

import (
	"errors"
)

type CreateOrderDTO struct {
	CustomerID string `json:"customer_id"`
	Products   []uint `json:"products"`
}

func (c *CreateOrderDTO) Validate() error {
	if c.CustomerID == "" {
		return errors.New("customer_id is required")
	}
	if len(c.Products) == 0 {
		return errors.New("at least one product is required")
	}
	for _, id := range c.Products {
		if id == 0 {
			return errors.New("product ID must be greater than zero")
		}
	}
	return nil
}
