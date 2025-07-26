package model

import (
	"testing"

	"github.com/fiap-161/tech-challenge-fiap161/hexagonal/adapters/drivers/rest/dto"
)

func TestAdmin_FromRegisterDTO(t *testing.T) {
	tests := []struct {
		name     string
		input    dto.RegisterDTO
		expected Admin
	}{
		{
			name: "basic register",
			input: dto.RegisterDTO{
				Email:    "admin@example.com",
				Password: "123456",
			},
			expected: Admin{
				Email:    "admin@example.com",
				Password: "123456",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Admin{}.FromRegisterDTO(tt.input)
			if result.Email != tt.expected.Email {
				t.Errorf("expected email %s, got %s", tt.expected.Email, result.Email)
			}
			if result.Password != tt.expected.Password {
				t.Errorf("expected password %s, got %s", tt.expected.Password, result.Password)
			}
		})
	}
}

func TestAdmin_FromLoginDTO(t *testing.T) {
	tests := []struct {
		name     string
		input    dto.LoginDTO
		expected Admin
	}{
		{
			name: "basic login",
			input: dto.LoginDTO{
				Email:    "user@example.com",
				Password: "abc123",
			},
			expected: Admin{
				Email:    "user@example.com",
				Password: "abc123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Admin{}.FromLoginDTO(tt.input)
			if result.Email != tt.expected.Email {
				t.Errorf("expected email %s, got %s", tt.expected.Email, result.Email)
			}
			if result.Password != tt.expected.Password {
				t.Errorf("expected password %s, got %s", tt.expected.Password, result.Password)
			}
		})
	}
}

func TestAdmin_Build(t *testing.T) {
	tests := []struct {
		name     string
		base     Admin
		password string
	}{
		{
			name: "build from admin",
			base: Admin{
				Email: "admin@build.com",
			},
			password: "buildPass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.base.Build(tt.password)

			if result.Email != tt.base.Email {
				t.Errorf("expected email %s, got %s", tt.base.Email, result.Email)
			}
			if result.Password != tt.password {
				t.Errorf("expected password %s, got %s", tt.password, result.Password)
			}
			if result.ID == "" {
				t.Error("expected non-empty ID")
			}
			if result.CreatedAt.IsZero() || result.UpdatedAt.IsZero() {
				t.Error("expected non-zero timestamps")
			}
		})
	}
}
