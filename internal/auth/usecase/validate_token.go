package usecase

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/auth/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/auth/gateway"
)

type ValidateTokenUseCase struct {
	tokenGateway gateway.TokenGateway
}

func NewValidateTokenUseCase(tokenGateway gateway.TokenGateway) *ValidateTokenUseCase {
	return &ValidateTokenUseCase{
		tokenGateway: tokenGateway,
	}
}

func (uc *ValidateTokenUseCase) Execute(tokenString string) (*entity.CustomClaims, error) {
	return uc.tokenGateway.ValidateToken(tokenString)
}
