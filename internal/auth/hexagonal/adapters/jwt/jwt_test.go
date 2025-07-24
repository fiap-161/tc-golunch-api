package auth

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	secretKey := "test-secret-key"
	expiryDuration := time.Hour
	service := NewJWTService(secretKey, expiryDuration)

	token, err := service.GenerateToken("user123", "admin", map[string]any{"role": "super-admin"})
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Expected token to be non-empty")
	}
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name          string
		secretKey     string
		validateKey   string
		userID        string
		userType      string
		claims        map[string]any
		expiry        time.Duration
		sleep         time.Duration
		expectError   bool
		errorContains string
	}{
		{
			name:        "Valid token",
			secretKey:   "test-secret-key",
			validateKey: "test-secret-key",
			userID:      "user123",
			userType:    "admin",
			claims:      map[string]any{"role": "super-admin"},
			expiry:      time.Hour,
			sleep:       0,
			expectError: false,
		},
		{
			name:          "Wrong secret key",
			secretKey:     "test-secret-key",
			validateKey:   "wrong-secret-key",
			userID:        "user123",
			userType:      "admin",
			claims:        nil,
			expiry:        time.Hour,
			sleep:         0,
			expectError:   true,
			errorContains: "signature",
		},
		{
			name:          "Expired token",
			secretKey:     "test-secret-key",
			validateKey:   "test-secret-key",
			userID:        "user123",
			userType:      "admin",
			claims:        nil,
			expiry:        time.Millisecond,
			sleep:         10 * time.Millisecond,
			expectError:   true,
			errorContains: "expired",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generateService := NewJWTService(tt.secretKey, tt.expiry)
			validateService := NewJWTService(tt.validateKey, tt.expiry)

			token, err := generateService.GenerateToken(tt.userID, tt.userType, tt.claims)
			if err != nil {
				t.Fatalf("Failed to generate token: %v", err)
			}

			if tt.sleep > 0 {
				time.Sleep(tt.sleep)
			}

			claims, err := validateService.ValidateToken(token)

			if tt.expectError && err == nil {
				t.Error("Expected an error but got nil")
				return
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}

			if tt.expectError && err != nil && tt.errorContains != "" {
				if err.Error() == "" || !contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error containing '%s', got: '%v'", tt.errorContains, err)
				}
				return
			}

			if !tt.expectError {
				if claims.UserID != tt.userID {
					t.Errorf("Expected UserID to be %s, got %s", tt.userID, claims.UserID)
				}

				if claims.UserType != tt.userType {
					t.Errorf("Expected UserType to be %s, got %s", tt.userType, claims.UserType)
				}

				if tt.claims != nil {
					for k, v := range tt.claims {
						if claims.Custom[k] != v {
							t.Errorf("Expected custom claim %s to be %v, got %v", k, v, claims.Custom[k])
						}
					}
				}
			}
		})
	}
}

func TestInvalidTokenFormat(t *testing.T) {
	service := NewJWTService("test-secret-key", time.Hour)

	invalidTokens := []struct {
		name  string
		token string
	}{
		{"Empty token", ""},
		{"Malformed token", "not.valid.jwt"},
		{"Incomplete token", "header.payload"},
	}

	for _, tt := range invalidTokens {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.ValidateToken(tt.token)
			if err == nil {
				t.Errorf("Expected error for invalid token format, got nil")
			}
		})
	}
}

func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && strings.Contains(s, substr)
}
