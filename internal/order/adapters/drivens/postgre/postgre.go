package postgre

import (
	"context"
	"gorm.io/gorm"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/core/model"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(DB *gorm.DB) *Repository {
	return &Repository{
		DB: DB,
	}
}

func (p *Repository) Create(ctx context.Context, order model.Order) (model.Order, error) {
	tx := p.DB.Create(&order)
	if tx.Error != nil {
		return model.Order{}, tx.Error
	}

	return order, nil
}

func (p *Repository) GetAll(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	
	if err := p.DB.Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}
