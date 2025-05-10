package ports

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/user/core/model"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (model.User, error)
}
