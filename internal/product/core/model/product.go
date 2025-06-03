package model

import (
	"time"

	"errors"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
	"github.com/google/uuid"
)

type Product struct {
	entity.Entity
	Name          string        `json:"name"`
	Price         float64       `json:"price" gorm:"type:decimal(10,2)"`
	Description   string        `json:"description" gorm:"type:text"`
	PreparingTime uint          `json:"preparing_time" gorm:"type:integer"`
	Category      enum.Category `json:"category"`
	ImageURL      string        `json:"image_url" gorm:"type:varchar(255)"`
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

func (p Product) Validate() error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Price < 0 {
		return errors.New("price must be positive")
	}
	if p.Category == 0 {
		return errors.New("category is required")
	}
	return nil
}
