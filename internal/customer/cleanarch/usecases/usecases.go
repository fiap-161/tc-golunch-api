package usecases

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/cleanarch/gateway"
	apperror "github.com/fiap-161/tech-challenge-fiap161/internal/shared/errors"
)

type CustomerUseCases struct {
	CustomerGateway gateway.CustomerGateway
}

func Build(gateway gateway.CustomerGateway) *CustomerUseCases {
	return &CustomerUseCases{
		CustomerGateway: gateway,
	}
}

func (u *CustomerUseCases) Create(ctx context.Context, customer entity.Customer) (string, error) {
	exists, _ := u.CustomerGateway.FindByCPF(ctx, customer.CPF)
	if exists.CPF != "" {
		return "", &apperror.ValidationError{Msg: "Customer already registered"}
	}

	customerWithID := customer.Build() // adiciona ID e timestamps
	saved, err := u.CustomerGateway.Create(ctx, customerWithID)
	if err != nil {
		return "", err
	}

	return saved.Id, nil
}

func (u *CustomerUseCases) FindByCPF(ctx context.Context, cpf string) (entity.Customer, error) {
	customer, err := u.CustomerGateway.FindByCPF(ctx, cpf)
	if err != nil || customer.Id == "" {
		return entity.Customer{}, &apperror.NotFoundError{Msg: "Customer not found"}
	}

	return customer, nil
}
