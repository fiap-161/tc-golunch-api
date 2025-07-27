package entity

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
)

type Order struct {
	entity.Entity
	CustomerID    string           `json:"customer_id" gorm:"index"`
	Status        enum.OrderStatus `json:"status" gorm:"type:varchar(20)"`
	Price         float64          `json:"price" gorm:"type:decimal(10,2)"`
	PreparingTime uint             `json:"preparing_time" gorm:"type:integer"`
}
