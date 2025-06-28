package entity

import (
	"strings"
	"time"

	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity/enum"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
	"github.com/google/uuid"
)

type Product struct {
	Id            string
	Name          string
	Price         float64
	Description   string
	PreparingTime uint
	Category      enum.Category
	ImageURL      string
}

func (p Product) Build() Product {
	return Product{
		Id:            p.Id,
		Name:          p.Name,
		Price:         p.Price,
		Description:   p.Description,
		PreparingTime: p.PreparingTime,
		Category:      p.Category,
		ImageURL:      p.ImageURL,
	}
}

func FromRequestDTO(dto dto.ProductRequestDTO) Product {
	category := strings.ToUpper(string(dto.Category))
	return Product{
		Name:          dto.Name,
		Price:         dto.Price,
		Description:   dto.Description,
		PreparingTime: dto.PreparingTime,
		Category:      enum.Category(category),
		ImageURL:      dto.ImageURL,
	}
}

func FromUpdateDTO(dto dto.ProductRequestUpdateDTO) Product {
	category := strings.ToUpper(string(dto.Category))
	return Product{
		Name:          dto.Name,
		Price:         dto.Price,
		Description:   dto.Description,
		PreparingTime: dto.PreparingTime,
		Category:      enum.Category(category),
		ImageURL:      dto.ImageURL,
	}
}

func (p Product) ToProductDAO() dto.ProductDAO {
	return dto.ProductDAO{
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

func FromProductDAO(dao dto.ProductDAO) Product {
	category := strings.ToUpper(string(dao.Category))
	return Product{
		Id:            dao.ID,
		Name:          dao.Name,
		Price:         dao.Price,
		Description:   dao.Description,
		PreparingTime: dao.PreparingTime,
		Category:      enum.Category(category),
		ImageURL:      dao.ImageURL,
	}
}

func (p Product) Validate() error {
	if p.Name == "" {
		return &apperror.ValidationError{Msg: "Name is required"}
	}
	if p.Price < 0 {
		return &apperror.ValidationError{Msg: "Price must be positive"}
	}
	if p.Category == "" {
		return &apperror.ValidationError{Msg: "Category is required"}
	}
	return nil
}
