package model

import (
	"time"

	"github.com/google/uuid"
	
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
)

type Admin struct {
	entity.Entity
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
}

func (a Admin) Build(password string) Admin {
	return Admin{
		Entity: entity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:    a.Email,
		Password: password,
	}
}

func (a Admin) FromRegisterDTO(dto dto.RegisterDTO) Admin {
	return Admin{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func (a Admin) FromLoginDTO(dto dto.LoginDTO) Admin {
	return Admin{
		Email:    dto.Email,
		Password: dto.Password,
	}
}
