package ports

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/core/model"
)

type AdminRepository interface {
	Create(ctx context.Context, admin model.Admin) error
	FindByEmail(ctx context.Context, email string) (model.Admin, error)
}
