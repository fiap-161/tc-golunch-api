package dto

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"
)

type ProductRequestDTO struct {
	Name          string  `json:"name" binding:"required"`
	Price         float64 `json:"price" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	PreparingTime uint    `json:"preparing_time" binding:"required"`
	Category      *uint   `json:"category" binding:"required"`
	ImageURL      string  `json:"image_url" binding:"required,url"`
}

type ProductResponseDTO struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Description   string  `json:"description"`
	PreparingTime uint    `json:"preparing_time"`
	Category      string  `json:"category"`
	ImageURL      string  `json:"image_url"`
}

func FromRequestDTOToModel(dto ProductRequestDTO) model.Product {
	return model.Product{
		Name:          dto.Name,
		Price:         dto.Price,
		Description:   dto.Description,
		PreparingTime: dto.PreparingTime,
		Category:      enum.Category(*dto.Category),
		ImageURL:      dto.ImageURL,
	}
}

func FromModelToResponseDTO(model model.Product) ProductResponseDTO {
	return ProductResponseDTO{
		ID:            model.ID,
		Name:          model.Name,
		Price:         model.Price,
		Description:   model.Description,
		PreparingTime: model.PreparingTime,
		Category:      model.Category.String(),
		ImageURL:      model.ImageURL,
	}
}
