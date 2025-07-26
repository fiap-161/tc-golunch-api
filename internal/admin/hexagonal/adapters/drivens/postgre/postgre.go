package postgre

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/hexagonal/core/ports"

	"gorm.io/gorm"

	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/hexagonal/core/model"
)

type DB interface {
	Create(value any) *gorm.DB
	Where(query any, args ...any) *gorm.DB
	First(dest any, conds ...any) *gorm.DB
}

type Repository struct {
	db DB
}

func NewRepository(db DB) ports.AdminRepository {
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
