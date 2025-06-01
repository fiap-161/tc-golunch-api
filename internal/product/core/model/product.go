package model

import (
	"time"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
	"github.com/google/uuid"
)

type Product struct {
	entity.Entity
	Name          string
	Price         float64
	Description   string
	PreparingTime uint
	Category      enum.Category
	ImageURL      string
}

func (p Product) Build() Product {
	return Product{
		Entity: entity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:          p.Name,
		Price:         p.Price,
		Description:   p.Description,
		PreparingTime: p.PreparingTime,
		Category:      p.Category,
		ImageURL:      p.ImageURL,
	}
}

func (p Product) FromRequestDTO(dto dto.ProductRequestDTO) Product {
	return Product{
		Name:          dto.Name,
		Price:         dto.Price,
		Description:   dto.Description,
		PreparingTime: dto.PreparingTime,
		Category:      enum.Category(dto.Category),
		ImageURL:      dto.ImageURL,
	}
}

func (p Product) FromUpdateDTO(dto dto.ProductRequestUpdateDTO) Product {
	return Product{
		Name:          dto.Name,
		Price:         dto.Price,
		Description:   dto.Description,
		PreparingTime: dto.PreparingTime,
		Category:      enum.Category(dto.Category),
		ImageURL:      dto.ImageURL,
	}
}

func (p Product) FromEntityToResponseDTO() dto.ProductResponseDTO {
	return dto.ProductResponseDTO{
		ID:            p.ID,
		Name:          p.Name,
		Price:         p.Price,
		Description:   p.Description,
		PreparingTime: p.PreparingTime,
		Category:      p.Category.String(),
		ImageURL:      p.ImageURL,
	}
}
