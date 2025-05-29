package dto

import (
	"gorm.io/gorm"
	
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"
)

type Product struct {
	gorm.Model
	Name          string
	Price         float64
	Description   string
	PreparingTime uint
	Category      string
	ImageURL      string
}

func FromDAOToModel(dao Product) model.Product {
	cat, _ := enum.FromCategoryString(dao.Category)
	return model.Product{
		ID:            dao.ID,
		Name:          dao.Name,
		Description:   dao.Description,
		ImageURL:      dao.ImageURL,
		Price:         dao.Price,
		PreparingTime: dao.PreparingTime,
		Category:      enum.Category(cat),
	}
}

func FromModelToDAO(model model.Product) Product {
	return Product{
		Name:          model.Name,
		Price:         model.Price,
		Description:   model.Description,
		Category:      model.Category.String(),
		ImageURL:      model.ImageURL,
		PreparingTime: model.PreparingTime,
	}
}
