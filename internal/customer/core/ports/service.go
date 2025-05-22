package ports

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest/dto"
)

type CustomerService interface {
	Create(ctx context.Context, customerDTO dto.CreateCustomerDTO) (string, error)
	Identify(ctx context.Context, CPF string) (string, error)
}
