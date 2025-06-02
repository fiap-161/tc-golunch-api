package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fiap-161/tech-challenge-fiap161/internal/auth/adapters/jwt"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/shared/entity"
)

type mockRepo struct {
	CreateFunc    func(ctx context.Context, customer model.Customer) (model.Customer, error)
	FindByCPFFunc func(ctx context.Context, cpf string) (model.Customer, error)
}

func (m *mockRepo) Create(ctx context.Context, customer model.Customer) (model.Customer, error) {
	return m.CreateFunc(ctx, customer)
}

func (m *mockRepo) FindByCPF(ctx context.Context, cpf string) (model.Customer, error) {
	return m.FindByCPFFunc(ctx, cpf)
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name      string
		input     dto.CreateCustomerDTO
		createErr error
		wantErr   bool
	}{
		{
			name: "successful create",
			input: dto.CreateCustomerDTO{
				Name: "John Doe",
				CPF:  "12345678900",
			},
			createErr: nil,
			wantErr:   false,
		},
		{
			name: "repo create error",
			input: dto.CreateCustomerDTO{
				Name: "Jane Doe",
				CPF:  "09876543211",
			},
			createErr: errors.New("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		repo := &mockRepo{
			CreateFunc: func(ctx context.Context, customer model.Customer) (model.Customer, error) {
				if tt.createErr != nil {
					return model.Customer{}, tt.createErr
				}
				customer.ID = "generated-id"
				return customer, nil
			},
		}

		service := New(repo, nil)
		id, err := service.Create(context.Background(), tt.input)

		if (err != nil) != tt.wantErr {
			t.Errorf("Create() %q failed: got err = %v, wantErr = %v", tt.name, err, tt.wantErr)
		}
		if !tt.wantErr && id == "" {
			t.Errorf("Create() %q failed: expected an ID but got empty string", tt.name)
		}
	}
}

func TestService_Identify(t *testing.T) {
	customer := model.Customer{
		Entity: entity.Entity{
			ID: "customer-id",
		},
		CPF: "12345678900",
	}

	tests := []struct {
		name        string
		inputCPF    string
		findErr     error
		foundCust   model.Customer
		expectErr   bool
		expectToken bool
	}{
		{
			name:        "identify with existing customer",
			inputCPF:    "12345678900",
			foundCust:   customer,
			expectErr:   false,
			expectToken: true,
		},
		{
			name:        "identify with unknown customer CPF",
			inputCPF:    "00000000000",
			findErr:     errors.New("not found"),
			expectErr:   true,
			expectToken: false,
		},
		{
			name:        "identify with empty CPF (anonymous)",
			inputCPF:    "",
			expectErr:   false,
			expectToken: true,
		},
	}

	for _, tt := range tests {
		repo := &mockRepo{
			FindByCPFFunc: func(ctx context.Context, cpf string) (model.Customer, error) {
				return tt.foundCust, tt.findErr
			},
		}

		jwtService := auth.NewJWTService("secret", time.Minute)
		service := New(repo, jwtService)

		token, err := service.Identify(context.Background(), tt.inputCPF)

		if (err != nil) != tt.expectErr {
			t.Errorf("Identify() %q: expected error = %v, got %v", tt.name, tt.expectErr, err)
		}
		if (token != "") != tt.expectToken {
			t.Errorf("Identify() %q: expected token presence = %v, got %v", tt.name, tt.expectToken, token != "")
		}
	}
}
