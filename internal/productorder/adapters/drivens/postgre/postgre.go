package postgre

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/productorder/core/model"
	"gorm.io/gorm"
)

type DB interface {
	Create(value any) *gorm.DB
	Where(query any, args ...any) *gorm.DB
	First(dest any, conds ...any) *gorm.DB
}

type Repository struct {
	db DB
}

func NewRepository(db DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateBulk(orders []model.ProductOrder) (int, error) {
	tx := r.db.Create(&orders)
	if tx.Error != nil {
		return 0, tx.Error
	}

	return len(orders), nil
}

func (r *Repository) FindByOrderID(orderID string) ([]model.ProductOrder, error) {
	var orders []model.ProductOrder

	tx := r.db.Where("order_id = ?", orderID).Find(&orders)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return orders, nil
}
