package postgre

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivens/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"
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
