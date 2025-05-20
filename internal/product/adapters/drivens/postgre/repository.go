package postgre

import (
	"errors"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivens/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
	appErrors "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) Create(product model.Product) (model.Product, error) {
	productDAO := dto.FromModelToDAO(product)
	result := r.DB.Create(&productDAO)
	if result.Error != nil {
		return model.Product{}, result.Error
	}
	return dto.FromDAOToModel(productDAO), nil
}

func (r *ProductRepository) GetAll(category string) ([]model.Product, error) {
	var productDAOs []dto.ProductDAO
	query := r.DB

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&productDAOs).Error; err != nil {
		return nil, err
	}

	var products []model.Product
	for _, dao := range productDAOs {
		products = append(products, dto.FromDAOToModel(dao))
	}

	return products, nil
}

func (r *ProductRepository) Update(id uint, updated model.Product) (model.Product, error) {
	var existing dto.ProductDAO
	if err := r.DB.First(&existing, id).Error; err != nil {
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
	if updated.Category.String() != "" {
		updates["category"] = updated.Category.String()
	}

	if len(updates) == 0 {
		return dto.FromDAOToModel(existing), nil
	}

	if err := r.DB.Model(&dto.ProductDAO{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return model.Product{}, err
	}

	var updatedDAO dto.ProductDAO
	if err := r.DB.First(&updatedDAO, id).Error; err != nil {
		return model.Product{}, err
	}

	return dto.FromDAOToModel(updatedDAO), nil
}

func (r *ProductRepository) FindByID(id uint) (model.Product, error) {
	var existing dto.ProductDAO
	if err := r.DB.First(&existing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Product{}, &appErrors.NotFoundError{Msg: "Product not found"}
		}
		return model.Product{}, err
	}

	return dto.FromDAOToModel(existing), nil
}

func (r *ProductRepository) Delete(id uint) error {
	var product dto.ProductDAO

	if err := r.DB.First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &appErrors.NotFoundError{Msg: "Product not found"}
		}
		return err
	}

	if err := r.DB.Delete(&product).Error; err != nil {
		return err
	}

	return nil

}
