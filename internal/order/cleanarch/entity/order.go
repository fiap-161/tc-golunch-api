package entity

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity/enum"
	productentity "github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
	"github.com/google/uuid"
	"time"
)

type Order struct {
	entity.Entity
	CustomerID    string           `json:"customer_id" gorm:"index"`
	Status        enum.OrderStatus `json:"status" gorm:"type:varchar(20)"`
	Price         float64          `json:"price" gorm:"type:decimal(10,2)"`
	PreparingTime uint             `json:"preparing_time" gorm:"type:integer"`
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

func (o Order) FromDTO(dto dto.CreateOrderDTO, products []productentity.Product) Order {
	totalPrice, preparingTime := o.getOrderInfoFromProducts(products, dto)

	return Order{
		CustomerID:    dto.CustomerID,
		Price:         totalPrice,
		PreparingTime: preparingTime,
		Status:        enum.OrderStatusAwaitingPayment,
	}
}

func (o Order) getOrderInfoFromProducts(products []productentity.Product, orderDTO dto.CreateOrderDTO) (float64, uint) {
	var totalPrice float64
	var preparingTime uint

	for _, item := range orderDTO.Products {
		for _, product := range products {
			if product.Id == item.ProductID {
				totalPrice += product.Price * float64(item.Quantity)
				preparingTime += product.PreparingTime
			}
		}
	}

	return totalPrice, preparingTime
}
