package dto

import "gorm.io/gorm"

type ProductDAO struct {
	gorm.Model
	Name          string
	Price         float64
	Description   string
	PreparingTime uint
	Category      string
	ImageURL      string
}
