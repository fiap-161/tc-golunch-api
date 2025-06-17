package model

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/hexagonal/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
)

type OrderStatus string

const (
	OrderStatusAwaitingPayment OrderStatus = "awaiting_payment"
	OrderStatusReceived        OrderStatus = "received"
	OrderStatusInPreparation   OrderStatus = "in_preparation"
	OrderStatusReady           OrderStatus = "ready"
	OrderStatusCompleted       OrderStatus = "completed"
)

var OrderPanelStatus = []string{
	OrderStatusReceived.String(),
	OrderStatusInPreparation.String(),
	OrderStatusReady.String(),
}

func (o OrderStatus) String() string {
	return string(o)
}

type Order struct {
	entity.Entity
	CustomerID    string      `json:"customer_id" gorm:"index"`
	Status        OrderStatus `json:"status" gorm:"type:varchar(20)"`
	Price         float64     `json:"price" gorm:"type:decimal(10,2)"`
	PreparingTime uint        `json:"preparing_time" gorm:"type:integer"`
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

func (o Order) BuildUpdate(status OrderStatus) Order {
	return Order{
		Entity: entity.Entity{
			ID:        o.ID,
			CreatedAt: o.CreatedAt,
			UpdatedAt: time.Now(),
		},
		CustomerID:    o.CustomerID,
		Status:        status,
		Price:         o.Price,
		PreparingTime: o.PreparingTime,
	}
}

func (o Order) FromDTO(dto dto.CreateOrderDTO, products []model.Product) Order {
	totalPrice, preparingTime := o.getOrderInfoFromProducts(products, dto)

	return Order{
		CustomerID:    dto.CustomerID,
		Price:         totalPrice,
		PreparingTime: preparingTime,
		Status:        OrderStatusAwaitingPayment,
	}
}

func (o Order) getOrderInfoFromProducts(products []model.Product, dto dto.CreateOrderDTO) (float64, uint) {
	var totalPrice float64
	var totalPreparingTime uint

	for _, item := range dto.Products {
		for _, product := range products {
			if product.ID == item.ProductID {
				totalPrice += product.Price * float64(item.Quantity)
				totalPreparingTime += product.PreparingTime * uint(item.Quantity)
			}
		}
	}

	return totalPrice, totalPreparingTime
}

func (o Order) ToPanelItemDTO() dto.OrderPanelItemDTO {
	return dto.OrderPanelItemDTO{
		OrderNumber:   strings.ToUpper(o.ID[len(o.ID)-4:]),
		Status:        o.Status.String(),
		PreparingTime: o.PreparingTime,
	}
}
