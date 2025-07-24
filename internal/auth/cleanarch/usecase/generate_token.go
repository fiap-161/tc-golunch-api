package usecase

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/auth/cleanarch/gateway"
)

type GenerateTokenUseCase struct {
	TokenGateway gateway.TokenGateway
}

func NewGenerateTokenUseCase(tokenGateway gateway.TokenGateway) *GenerateTokenUseCase {
	return &GenerateTokenUseCase{TokenGateway: tokenGateway}
}

func (uc *GenerateTokenUseCase) Execute(userID, userType string, additionalClaims map[string]any) (string, error) {
	return uc.TokenGateway.GenerateToken(userID, userType, additionalClaims)
} 