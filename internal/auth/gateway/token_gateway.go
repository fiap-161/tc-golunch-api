package gateway

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/auth/entity"
)

type TokenGateway interface {
	GenerateToken(userID, userType string, additionalClaims map[string]any) (string, error)
	ValidateToken(tokenString string) (*entity.CustomClaims, error)
}
