package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID   string         `json:"user_id"`
	UserType string         `json:"user_type"`
	Custom   map[string]any `json:"custom"`
}

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

	claims := CustomClaims{
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

func (s *JWTService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
