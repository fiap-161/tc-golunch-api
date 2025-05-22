package entity

import (
	"time"
)

type Entity struct {
	ID        string    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
