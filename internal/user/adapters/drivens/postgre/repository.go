package postgre

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/user/core/model"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

func (p *repository) GetUserByID(_ context.Context, id string) (model.User, error) {
	var user model.User

	tx := p.db.Where("id = ?", id).First(&user)

	if tx.Error != nil {
		return model.User{}, tx.Error
	}

	return user, nil
}

func (p *repository) Create(_ context.Context, user model.User) error {
	tx := p.db.Create(&user)
	if tx.Error != nil {
		return tx.Error
	}
	
	return nil
}
