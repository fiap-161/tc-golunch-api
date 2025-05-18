package model

import (
	"testing"
	"time"

	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest/dto"
)

func TestCustomer_Build(t *testing.T) {
	tests := []struct {
		name            string
		input           Customer
		wantName        string
		wantEmail       string
		wantCPF         string
		wantIsAnonymous bool
	}{
		{
			name: "build with name, email, cpf, anonymous false",
			input: Customer{
				Name:        "John Doe",
				Email:       "john@example.com",
				CPF:         "12345678900",
				IsAnonymous: false,
			},
			wantName:        "John Doe",
			wantEmail:       "john@example.com",
			wantCPF:         "12345678900",
			wantIsAnonymous: false,
		},
		{
			name: "build with anonymous true",
			input: Customer{
				Name:        "Anon",
				Email:       "anon@example.com",
				CPF:         "00000000000",
				IsAnonymous: true,
			},
			wantName:        "Anon",
			wantEmail:       "anon@example.com",
			wantCPF:         "00000000000",
			wantIsAnonymous: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Build()

			if got.ID == "" {
				t.Errorf("Build() ID = empty, want non-empty")
			}

			now := time.Now()
			if got.CreatedAt.Before(now.Add(-1*time.Minute)) || got.CreatedAt.After(now.Add(1*time.Minute)) {
				t.Errorf("Build() CreatedAt = %v, want close to now", got.CreatedAt)
			}
			if got.UpdatedAt.Before(now.Add(-1*time.Minute)) || got.UpdatedAt.After(now.Add(1*time.Minute)) {
				t.Errorf("Build() UpdatedAt = %v, want close to now", got.UpdatedAt)
			}

			if got.Name != tt.wantName {
				t.Errorf("Build() Name = %v, want %v", got.Name, tt.wantName)
			}
			if got.Email != tt.wantEmail {
				t.Errorf("Build() Email = %v, want %v", got.Email, tt.wantEmail)
			}
			if got.CPF != tt.wantCPF {
				t.Errorf("Build() CPF = %v, want %v", got.CPF, tt.wantCPF)
			}
			if got.IsAnonymous != tt.wantIsAnonymous {
				t.Errorf("Build() IsAnonymous = %v, want %v", got.IsAnonymous, tt.wantIsAnonymous)
			}
		})
	}
}

func TestCustomer_FromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  dto.CreateCustomerDTO
		want Customer
	}{
		{
			name: "from dto normal",
			dto: dto.CreateCustomerDTO{
				Name:  "Alice",
				Email: "alice@example.com",
				CPF:   "11122233344",
			},
			want: Customer{
				Name:  "Alice",
				Email: "alice@example.com",
				CPF:   "11122233344",
			},
		},
		{
			name: "from dto empty",
			dto: dto.CreateCustomerDTO{
				Name:  "",
				Email: "",
				CPF:   "",
			},
			want: Customer{
				Name:  "",
				Email: "",
				CPF:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Customer{}.FromDTO(tt.dto)

			if got.Name != tt.want.Name {
				t.Errorf("FromDTO() Name = %v, want %v", got.Name, tt.want.Name)
			}
			if got.Email != tt.want.Email {
				t.Errorf("FromDTO() Email = %v, want %v", got.Email, tt.want.Email)
			}
			if got.CPF != tt.want.CPF {
				t.Errorf("FromDTO() CPF = %v, want %v", got.CPF, tt.want.CPF)
			}
		})
	}
}
