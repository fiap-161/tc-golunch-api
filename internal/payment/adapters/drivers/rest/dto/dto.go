package dto

type CheckPaymentDTO struct {
	Resource string `json:"resource" binding:"required"`
	Topic    string `json:"topic" binding:"required"`
}
