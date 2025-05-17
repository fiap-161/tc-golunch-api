package postgre

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/core/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(_ context.Context, admin model.Admin) error {
	tx := r.db.Create(&admin)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *Repository) FindByEmail(_ context.Context, email string) (model.Admin, error) {
	var admin model.Admin

	tx := r.db.Where("email = ?", email).First(&admin)

	if tx.Error != nil {
		return model.Admin{}, tx.Error
	}

	return admin, nil
}
