package model

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
	"time"

	"github.com/google/uuid"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
)

type OrderStatus string

const (
	OrderStatusReceived      OrderStatus = "received"
	OrderStatusInPreparation OrderStatus = "in_preparation"
	OrderStatusReady         OrderStatus = "ready"
	OrderStatusCompleted     OrderStatus = "completed"
)

type Order struct {
	entity.Entity
	CustomerID    string      `json:"customer_id" gorm:"index"`
	Status        OrderStatus `json:"status" gorm:"type:varchar(20)"`
	Price         float64     `json:"price" gorm:"type:decimal(10,2)"`
	PreparingTime uint        `json:"preparing_time"`
}

func (o Order) Build() Order {
	return Order{
		Entity: entity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CustomerID:    o.CustomerID,
		Status:        o.Status,
		Price:         o.Price,
		PreparingTime: o.PreparingTime,
	}
}

func (o Order) FromDTO(dto dto.CreateOrderDTO, products []model.Product) Order {
	totalPrice, preparingTime := o.getOrderInfoFromProducts(products)

	return Order{
		CustomerID:    dto.CustomerID,
		Price:         totalPrice,
		PreparingTime: preparingTime,
		Status:        OrderStatusReceived,
	}
}

func (o Order) getOrderInfoFromProducts(products []model.Product) (float64, uint) {
	var totalPrice float64
	var totalPreparingTime uint

	for _, product := range products {
		totalPrice += product.Price
		totalPreparingTime += product.PreparingTime
	}

	return totalPrice, totalPreparingTime
}
