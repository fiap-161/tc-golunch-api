package ports

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/auth/core/model"
)

type TokenService interface {
	GenerateToken(userID, userType string, additionalClaims map[string]any) (string, error)
	ValidateToken(tokenString string) (*model.CustomClaims, error)
}
