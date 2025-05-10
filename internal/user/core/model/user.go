package model

import "time"

type User struct {
	ID        uint `gorm:"primary_key"`
	Name      string
	Email     string
	CPF       string
	CreatedAt time.Time
}
