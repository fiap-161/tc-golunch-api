package postgre

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/model"
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

func (p *Repository) FindByCPF(_ context.Context, CPF string) (model.Customer, error) {
	var customer model.Customer

	tx := p.db.Where("cpf = ?", CPF).First(&customer)

	if tx.Error != nil {
		return model.Customer{}, tx.Error
	}

	return customer, nil
}

func (p *Repository) Create(_ context.Context, customer model.Customer) (model.Customer, error) {
	tx := p.db.Create(&customer)
	if tx.Error != nil {
		return model.Customer{}, tx.Error
	}

	return customer, nil
}
