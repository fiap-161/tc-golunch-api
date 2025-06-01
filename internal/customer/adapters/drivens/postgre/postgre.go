package postgre

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/ports"

	"gorm.io/gorm"

	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/model"
)

type DB interface {
	Create(value any) *gorm.DB
	Where(query any, args ...any) *gorm.DB
	First(dest any, conds ...any) *gorm.DB
}

type Repository struct {
	db DB
}

func NewRepository(db DB) ports.CustomerRepository {
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
