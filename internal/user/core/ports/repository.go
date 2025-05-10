package ports

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/user/core/model"
)

type UserRepository interface {
	GetUserByID(_ context.Context, id string) (model.User, error)
}
