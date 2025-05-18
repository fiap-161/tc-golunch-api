package postgre

import (
	"errors"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivens/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"
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
	productDAO := dto.ProductDAO{
		Name:          product.Name,
		Price:         product.Price,
		Description:   product.Description,
		Category:      product.Category.String(),
		ImageURL:      product.ImageURL,
		PreparingTime: product.PreparingTime,
	}
	result := r.DB.Create(&productDAO)
	if result.Error != nil {
		return model.Product{}, result.Error
	}
	cat, _ := enum.FromCategoryString(productDAO.Category)
	return model.Product{
		ID:            productDAO.ID,
		Name:          productDAO.Name,
		Price:         productDAO.Price,
		Description:   productDAO.Description,
		PreparingTime: productDAO.PreparingTime,
		Category:      enum.Category(cat),
		ImageURL:      productDAO.ImageURL,
	}, nil
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
		cat, _ := enum.FromCategoryString(dao.Category)
		product := model.Product{
			ID:            dao.ID,
			Name:          dao.Name,
			Price:         dao.Price,
			Description:   dao.Description,
			ImageURL:      dao.ImageURL,
			PreparingTime: dao.PreparingTime,
			Category:      enum.Category(cat),
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) Update(id uint, updated model.Product) (model.Product, error) {
	var existing dto.ProductDAO
	if err := r.DB.First(&existing, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Product{}, &appErrors.NotFoundError{Msg: "Product not found"}
		}
		return model.Product{}, err
	}

	updates := map[string]interface{}{}
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
		cat, _ := enum.FromCategoryString(existing.Category)
		return model.Product{
			ID:            existing.ID,
			Name:          existing.Name,
			Description:   existing.Description,
			ImageURL:      existing.ImageURL,
			Price:         existing.Price,
			PreparingTime: existing.PreparingTime,
			Category:      enum.Category(cat),
		}, nil
	}

	if err := r.DB.Model(&dto.ProductDAO{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return model.Product{}, err
	}

	var updatedDAO dto.ProductDAO
	if err := r.DB.First(&updatedDAO, id).Error; err != nil {
		return model.Product{}, err
	}

	cat, _ := enum.FromCategoryString(updatedDAO.Category)
	return model.Product{
		ID:            updatedDAO.ID,
		Name:          updatedDAO.Name,
		Description:   updatedDAO.Description,
		ImageURL:      updatedDAO.ImageURL,
		Price:         updatedDAO.Price,
		PreparingTime: updatedDAO.PreparingTime,
		Category:      enum.Category(cat),
	}, nil

}
