package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/adapters/drivers/rest/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/model"
	"github.com/fiap-161/tech-challenge-fiap161/internal/customer/core/ports"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type Service struct {
	repo ports.CustomerRepository
}

func New(repo ports.CustomerRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, customerDTO dto.CreateCustomerDTO) (string, error) {
	var customer model.Customer
	customer = customer.FromDTO(customerDTO)

	createdCustomer, err := s.repo.Create(ctx, customer.Build())
	if err != nil {
		return "", err
	}

	return createdCustomer.Entity.ID, nil
}

func (s *Service) Identify(ctx context.Context, CPF string) (string, error) {
	customer, err := s.repo.FindByCPF(ctx, CPF)
	if err != nil {
		return "", errors.New("customer not found")
	}

	token, err := createToken(customer.Entity.ID, false)
	if err != nil {
		return "", errors.New("error creating token")
	}

	return token, nil
}

var secretKey = []byte(os.Getenv("SECRET_KEY"))

func createToken(id string, isAnonymous bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":           id,
			"is_anonymous": isAnonymous,
			"exp":          time.Now().Add(time.Hour * 24).Unix(),
		})

	fmt.Print()
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
