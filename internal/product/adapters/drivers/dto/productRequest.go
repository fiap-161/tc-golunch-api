package dto

type ProductRequestDTO struct {
	Name          string  `json:"name" binding:"required"`
	Price         float64 `json:"price" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	PreparingTime uint    `json:"preparing_time" binding:"required"`
	Category      *uint   `json:"category" binding:"required"`
	ImageURL      string  `json:"image_url" binding:"required,url"`
}
