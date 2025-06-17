package presenter

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/product/cleanarch/entity"
)

type Presenter struct {
}

func Build() *Presenter {
	return &Presenter{}
}

func (p *Presenter) FromEntityToResponseDTO(product entity.Product) dto.ProductResponseDTO {
	return dto.ProductResponseDTO{
		ID:            product.ID,
		Name:          product.Name,
		Price:         product.Price,
		Description:   product.Description,
		PreparingTime: product.PreparingTime,
		Category:      product.Category,
		ImageURL:      product.ImageURL,
	}
}
