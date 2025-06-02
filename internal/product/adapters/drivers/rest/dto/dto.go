package dto

type ProductRequestDTO struct {
	Name          string  `json:"name" binding:"required"`
	Price         float64 `json:"price" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	PreparingTime uint    `json:"preparing_time" binding:"required"`
	Category      uint    `json:"category_id" binding:"required"`
	ImageURL      string  `json:"image_url" binding:"required,url"`
}

type ProductResponseDTO struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Description   string  `json:"description"`
	PreparingTime uint    `json:"preparing_time"`
	Category      string  `json:"category"`
	ImageURL      string  `json:"image_url"`
}

type ProductListResponseDTO struct {
	Total uint                 `json:"total"`
	List  []ProductResponseDTO `json:"list"`
}

type ProductRequestUpdateDTO struct {
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Description   string  `json:"description"`
	PreparingTime uint    `json:"preparing_time"`
	Category      uint    `json:"category_id"`
	ImageURL      string  `json:"image_url" binding:"url"`
}
