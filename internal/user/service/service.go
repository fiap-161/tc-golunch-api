package service

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/user/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/user/core/ports"
)

type Service struct {
	repo ports.UserRepository
}

func New(repo ports.UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetUserByID(ctx context.Context, id string) (model.User, error) {
	return s.repo.GetUserByID(ctx, id)
}
