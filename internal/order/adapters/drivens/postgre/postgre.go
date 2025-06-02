package postgre

import (
	"context"

	"gorm.io/gorm"

	"github.com/fiap-161/tech-challenge-fiap161/internal/order/core/model"
	orderport "github.com/fiap-161/tech-challenge-fiap161/internal/order/core/ports"
)

type Repository struct {
	DB *gorm.DB
}

func New(DB *gorm.DB) orderport.OrderRepository {
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

func (p *Repository) FindByID(ctx context.Context, id string) (model.Order, error) {
	var order model.Order

	tx := p.DB.First(&order, "id = ?", id)
	if tx.Error != nil {
		return model.Order{}, tx.Error
	}

	return order, nil
}

func (p *Repository) GetPanel(ctx context.Context, status []string) ([]model.Order, error) {
	var orders []model.Order

	if err := p.DB.Where("status IN ?", status).Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (p *Repository) Update(ctx context.Context, order model.Order) (model.Order, error) {
	tx := p.DB.Save(&order)
	if tx.Error != nil {
		return model.Order{}, tx.Error
	}

	return order, nil
}
