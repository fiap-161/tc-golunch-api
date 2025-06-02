package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
)

type Payment struct {
	entity.Entity
	OrderID string `json:"order_id" gorm:"not null;unique"`
	QrCode  string `json:"qr_code" gorm:"not null"`
	Status  Status `json:"status" gorm:"not null;default:'PENDING'"`
}

type Status string

const (
	Pending  Status = "PENDING"
	Approved Status = "APPROVED"
	Rejected Status = "REJECTED"
)

func (p Payment) Build(orderID, qrCode string) Payment {
	now := time.Now()

	return Payment{
		Entity: entity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		OrderID: orderID,
		QrCode:  qrCode,
		Status:  Pending,
	}
}
