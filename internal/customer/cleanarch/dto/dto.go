package dto

import (
	"time"

	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/cleanarch/entity"
)

type CustomerDAO struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CPF       string    `gorm:"uniqueIndex" json:"cpf"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToCustomerDAO(c entity.Customer) CustomerDAO {
	return CustomerDAO{
		ID:        c.Id,
		Name:      c.Name,
		Email:     c.Email,
		CPF:       c.CPF,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func FromCustomerDAO(d CustomerDAO) entity.Customer {
	return entity.Customer{
		Id:    d.ID,
		Name:  d.Name,
		Email: d.Email,
		CPF:   d.CPF,
	}
}

type CreateCustomerDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	CPF   string `json:"cpf"`
}

type TokenDTO struct {
	TokenString string `json:"token"`
}
