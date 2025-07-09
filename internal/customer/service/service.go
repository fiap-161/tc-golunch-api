package service

import (
	"context"
	auth "github.com/fiap-161/tech-challenge-fiap161/internal/auth/hexagonal/core/ports"

	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/ports"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"

	"github.com/google/uuid"
)

type Service struct {
	repo       ports.CustomerRepository
	jwtService auth.TokenService
}

func New(repo ports.CustomerRepository, jwtService auth.TokenService) *Service {
	return &Service{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *Service) Create(ctx context.Context, customerDTO dto.CreateCustomerDTO) (string, error) {
	var customer model.Customer
	customer = customer.FromDTO(customerDTO)

	createdCustomer, err := s.repo.Create(ctx, customer.Build())
	if err != nil {
		return "", err
	}

	return createdCustomer.ID, nil
}

func (s *Service) Identify(ctx context.Context, CPF string) (string, error) {
	if CPF == "" {
		return s.createAnonymousToken()
	}

	customer, err := s.repo.FindByCPF(ctx, CPF)
	if err != nil {
		return "", &apperror.NotFoundError{Msg: "Customer not found"}
	}

	token, err := s.createToken(customer.ID, false)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) createAnonymousToken() (string, error) {
	anonymousID := uuid.NewString()

	token, err := s.createToken(anonymousID, true)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) createToken(id string, isAnonymous bool) (string, error) {
	additionalClaims := map[string]any{
		"is_anonymous": isAnonymous,
	}

	token, err := s.jwtService.GenerateToken(id, "customer", additionalClaims)
	if err != nil {
		return "", &apperror.InternalError{Msg: "Error creating token"}
	}

	return token, nil
}
