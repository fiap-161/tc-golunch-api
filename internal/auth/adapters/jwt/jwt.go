package auth

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/auth/core/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey      string
	expiryDuration time.Duration
}

func NewJWTService(secretKey string, expiryDuration time.Duration) *JWTService {
	return &JWTService{
		secretKey:      secretKey,
		expiryDuration: expiryDuration,
	}
}

func (s *JWTService) GenerateToken(userID, userType string, additionalClaims map[string]any) (string, error) {
	now := time.Now()

	claims := model.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID:   userID,
		UserType: userType,
		Custom:   additionalClaims,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (*model.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
