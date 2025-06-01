package postgre

import (
	"context"

	"gorm.io/gorm"

	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/payment/core/ports"
)

type DB interface {
	Create(value any) *gorm.DB
	Save(value any) *gorm.DB
	Where(query any, args ...any) *gorm.DB
	First(dest any, conds ...any) *gorm.DB
	Find(dest any, conds ...any) *gorm.DB
	Delete(value any, conds ...any) *gorm.DB
	Model(value any) *gorm.DB
	Updates(values any) *gorm.DB
}

type Repository struct {
	db DB
}

func New(db DB) ports.PaymentRepository {
	return &Repository{
		db: db,
	}
}

func (p *Repository) Create(ctx context.Context, payment model.Payment) (model.Payment, error) {
	tx := p.db.Create(&payment)
	if tx.Error != nil {
		return model.Payment{}, tx.Error
	}

	return payment, nil
}

func (p *Repository) GetAll(ctx context.Context) ([]model.Payment, error) {
	var payments []model.Payment

	if err := p.db.Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}

func (p *Repository) FindByOrderID(ctx context.Context, orderID string) (model.Payment, error) {
	var payment model.Payment

	tx := p.db.First(&payment, "order_id = ?", orderID)
	if tx.Error != nil {
		return model.Payment{}, tx.Error
	}

	return payment, nil
}

func (p *Repository) Update(ctx context.Context, payment model.Payment) (model.Payment, error) {
	tx := p.db.Save(&payment)
	if tx.Error != nil {
		return model.Payment{}, tx.Error
	}

	return payment, nil
}
