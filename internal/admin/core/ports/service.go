package ports

import (
	"context"
	
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivers/rest/dto"
)

type AdminService interface {
	Register(ctx context.Context, input dto.RegisterDTO) error
	Login(ctx context.Context, input dto.LoginDTO) (string, error) // returns JWT token
}
