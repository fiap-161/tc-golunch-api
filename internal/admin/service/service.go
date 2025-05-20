package service

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/core/ports"
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/utils"
	auth "github.com/fiap-161/tech-challenge-fiap161/internal/auth/core/ports"
	appErrors "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type Service struct {
	repo       ports.AdminRepository
	jwtService auth.TokenService
}

func New(repo ports.AdminRepository, jwtService auth.TokenService) *Service {
	return &Service{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *Service) Register(ctx context.Context, input dto.RegisterDTO) error {
	var admin model.Admin
	admin = admin.FromRegisterDTO(input)

	hashed, err := utils.HashPassword(admin.Password)
	if err != nil {
		return err
	}

	return s.repo.Create(ctx, admin.Build(hashed))
}

func (s *Service) Login(ctx context.Context, input dto.LoginDTO) (string, error) {
	var admin model.Admin
	admin = admin.FromLoginDTO(input)

	saved, err := s.repo.FindByEmail(ctx, admin.Email)
	if err != nil {
		return "", &appErrors.UnauthorizedError{Msg: "Invalid email or password"}
	}

	if !utils.CheckPasswordHash(input.Password, saved.Password) {
		return "", &appErrors.UnauthorizedError{Msg: "Invalid email or password"}
	}

	return s.jwtService.GenerateToken(saved.ID, "admin", nil)
}
