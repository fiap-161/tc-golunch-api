package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
)

type Customer struct {
	entity.Entity
	Name        string `json:"name"`
	Email       string `json:"email" gorm:"unique"`
	CPF         string `json:"cpf" gorm:"unique"`
	IsAnonymous bool   `json:"is_anonymous" gorm:"default:false"`
}

func (c Customer) Build() Customer {
	return Customer{
		Entity: entity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        c.Name,
		Email:       c.Email,
		CPF:         c.CPF,
		IsAnonymous: c.IsAnonymous,
	}
}

func (c Customer) FromDTO(dto dto.CreateCustomerDTO) Customer {
	return Customer{
		Name:  dto.Name,
		Email: dto.Email,
		CPF:   dto.CPF,
	}
}
