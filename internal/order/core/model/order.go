package model

import (
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

func (o Order) Build(price float64, preparingTime uint, productsJSON []byte) Order {
	return Order{
		Entity: entity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CustomerID:    o.CustomerID,
		Status:        OrderStatusReceived,
		Price:         price,
		PreparingTime: preparingTime,
	}
}

func (o Order) FromDTO(dto dto.CreateOrderDTO) Order {
	return Order{
		CustomerID: dto.CustomerID,
	}
}
