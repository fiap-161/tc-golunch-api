package dto

import (
	"errors"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/order/cleanarch/entity/enum"
	gormEntity "github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
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

type OrderDAO struct {
	gormEntity.Entity
	CustomerID    string           `json:"customer_id" gorm:"index"`
	Status        enum.OrderStatus `json:"status" gorm:"type:varchar(20)"`
	Price         float64          `json:"price" gorm:"type:decimal(10,2)"`
	PreparingTime uint             `json:"preparing_time" gorm:"type:integer"`
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

func ToOrderDAO(order entity.Order) OrderDAO {
	return OrderDAO{
		Entity:        order.Entity,
		CustomerID:    order.CustomerID,
		Status:        order.Status,
		Price:         order.Price,
		PreparingTime: order.PreparingTime,
	}
}

func FromOrderDAO(dao OrderDAO) entity.Order {
	return entity.Order{
		Entity:        dao.Entity,
		CustomerID:    dao.CustomerID,
		Status:        dao.Status,
		Price:         dao.Price,
		PreparingTime: dao.PreparingTime,
	}
}

func FromCreateOrderDTO(dto CreateOrderDTO) entity.Order {
	return entity.Order{
		CustomerID: dto.CustomerID,
		// Status, Price, PreparingTime podem ser definidos em outro lugar se necessário
	}
}

func EntityListFromDAOList(daoList []OrderDAO) []entity.Order {
	orders := make([]entity.Order, 0, len(daoList))
	for _, dao := range daoList {
		orders = append(orders, FromOrderDAO(dao))
	}
	return orders
}
