package service

import (
	"context"
	"errors"
	"testing"
	"time"
	
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/admin/utils"
	"github.com/fiap-161/tech-challenge-fiap161/internal/auth"
)

type mockRepo struct {
	CreateFunc      func(ctx context.Context, admin model.Admin) error
	FindByEmailFunc func(ctx context.Context, email string) (model.Admin, error)
}

func (m *mockRepo) Create(ctx context.Context, admin model.Admin) error {
	return m.CreateFunc(ctx, admin)
}

func (m *mockRepo) FindByEmail(ctx context.Context, email string) (model.Admin, error) {
	return m.FindByEmailFunc(ctx, email)
}

func TestService_Register(t *testing.T) {
	tests := []struct {
		name      string
		input     dto.RegisterDTO
		createErr error
		wantErr   bool
	}{
		{
			name: "successfully register",
			input: dto.RegisterDTO{
				Email:    "alice@example.com",
				Password: "password123",
			},
			createErr: nil,
			wantErr:   false,
		},
		{
			name: "repo create error",
			input: dto.RegisterDTO{
				Email:    "bob@example.com",
				Password: "password123",
			},
			createErr: errors.New("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		repo := &mockRepo{
			CreateFunc: func(ctx context.Context, admin model.Admin) error {
				return tt.createErr
			},
		}
		service := New(repo, nil)

		err := service.Register(context.Background(), tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("Register() %q failed: got err = %v, wantErr = %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestService_Login(t *testing.T) {
	hashedPwd, _ := utils.HashPassword("correctpassword")
	admin := model.Admin{
		Email:    "test@example.com",
		Password: hashedPwd,
	}

	tests := []struct {
		name        string
		input       dto.LoginDTO
		findErr     error
		foundAdmin  model.Admin
		expectErr   bool
		expectToken bool
	}{
		{
			name: "successful login",
			input: dto.LoginDTO{
				Email:    "test@example.com",
				Password: "correctpassword",
			},
			foundAdmin:  admin,
			expectErr:   false,
			expectToken: true,
		},
		{
			name: "invalid password",
			input: dto.LoginDTO{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			foundAdmin:  admin,
			expectErr:   true,
			expectToken: false,
		},
		{
			name: "admin not found",
			input: dto.LoginDTO{
				Email:    "notfound@example.com",
				Password: "any",
			},
			findErr:     errors.New("not found"),
			expectErr:   true,
			expectToken: false,
		},
	}

	for _, tt := range tests {
		repo := &mockRepo{
			FindByEmailFunc: func(ctx context.Context, email string) (model.Admin, error) {
				return tt.foundAdmin, tt.findErr
			},
		}

		jwtService := auth.NewJWTService("secret", time.Minute)
		service := New(repo, jwtService)

		token, err := service.Login(context.Background(), tt.input)
		if (err != nil) != tt.expectErr {
			t.Errorf("Login() %q: expected error = %v, got %v", tt.name, tt.expectErr, err)
		}
		if (token != "") != tt.expectToken {
			t.Errorf("Login() %q: expected token = %v, got %v", tt.name, tt.expectToken, token != "")
		}
	}
}
