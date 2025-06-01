package postgre

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type DB interface {
	Create(value any) *gorm.DB
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

func New(db DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(_ context.Context, product model.Product) (model.Product, error) {
	tx := r.db.Create(&product)
	if tx.Error != nil {
		return model.Product{}, tx.Error
	}

	return product, nil
}

func (r *Repository) GetAll(_ context.Context, category string) ([]model.Product, error) {
	var products []model.Product

	query := r.db
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repository) Update(_ context.Context, id string, updated model.Product) (model.Product, error) {
	var existing model.Product
	if err := r.db.First(&existing, id).Error; err != nil {
		return model.Product{}, err
	}

	updates := map[string]any{}
	if updated.Name != "" {
		updates["name"] = updated.Name
	}
	if updated.Description != "" {
		updates["description"] = updated.Description
	}
	if updated.ImageURL != "" {
		updates["image_url"] = updated.ImageURL
	}
	if updated.Price != 0 {
		updates["price"] = updated.Price
	}
	if updated.PreparingTime != 0 {
		updates["preparing_time"] = updated.PreparingTime
	}
	if updated.Category != enum.Unknown {
		updates["category"] = updated.Category.String()
	}

	if len(updates) == 0 {
		return existing, nil
	}

	if err := r.db.Model(&model.Product{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return model.Product{}, err
	}

	var updatedProduct model.Product
	if err := r.db.First(&updatedProduct, id).Error; err != nil {
		return model.Product{}, err
	}

	return updatedProduct, nil
}

func (r *Repository) FindByID(_ context.Context, id string) (model.Product, error) {
	var product model.Product
	if err := r.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Product{}, &apperror.NotFoundError{Msg: "Product not found"}
		}
		return model.Product{}, err
	}

	return product, nil
}

func (r *Repository) FindByIDs(_ context.Context, ids []string) ([]model.Product, error) {
	var products []model.Product

	if err := r.db.Where("id IN ?", ids).Find(&products).Error; err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, &apperror.NotFoundError{Msg: "No products found"}
	}

	return products, nil
}

func (r *Repository) Delete(_ context.Context, id string) error {
	var product model.Product

	if err := r.db.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &apperror.NotFoundError{Msg: "Product not found"}
		}
		return err
	}

	if err := r.db.Delete(&product).Error; err != nil {
		return err
	}

	return nil
}
