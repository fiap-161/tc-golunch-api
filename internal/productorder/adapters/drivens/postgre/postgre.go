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
