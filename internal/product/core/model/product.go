package model

import "github.com/fiap-161/tech-challenge-fiap161/internal/product/core/model/enum"

type Product struct {
	ID            uint
	Name          string
	Price         float64
	Description   string
	PreparingTime uint
	Category      enum.Category
	ImageURL      string
}
