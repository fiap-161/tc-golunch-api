package entity

import (
	"time"

	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
	"github.com/google/uuid"
)

type ProductOrder struct {
	ID        string
	ProductID string
	OrderID   string
	Quantity  int
	UnitPrice float64
}

func (po ProductOrder) ToProductOrderDAO() dto.ProductOrderDAO {
	return dto.ProductOrderDAO{
		Entity: entity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ProductID: po.ProductID,
		OrderID:   po.OrderID,
		Quantity:  po.Quantity,
		UnitPrice: po.UnitPrice,
	}
}

func ToListProducOrderDAO(list []ProductOrder) []dto.ProductOrderDAO {
	var listProductOrderDTO []dto.ProductOrderDAO
	for _, item := range list {
		dao := item.ToProductOrderDAO()
		listProductOrderDTO = append(listProductOrderDTO, dao)
	}
	return listProductOrderDTO
}

func FromRequestDTO(dto dto.ProductOrderRequestDTO) ProductOrder {
	return ProductOrder{
		ProductID: dto.ProductID,
		OrderID:   dto.OrderID,
		Quantity:  dto.Quantity,
		UnitPrice: dto.UnitPrice,
	}
}

func FromProductDAO(dao dto.ProductOrderDAO) ProductOrder {
	return ProductOrder{
		ID:        dao.ID,
		OrderID:   dao.OrderID,
		ProductID: dao.ProductID,
		Quantity:  dao.Quantity,
		UnitPrice: dao.UnitPrice,
	}
}

func ToListProducOrder(list []dto.ProductOrderDAO) []ProductOrder {
	var listProductOrder []ProductOrder
	for _, dao := range list {
		item := FromProductDAO(dao)
		listProductOrder = append(listProductOrder, item)
	}
	return listProductOrder
}
