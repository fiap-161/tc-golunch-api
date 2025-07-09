package usecases

import "github.com/fiap-161/tech-challenge-fiap161/internal/productorder/cleanarch/gateway"

type UseCases struct {
	ProductOrderGateway gateway.Gateway
}

func Build(productGateway gateway.Gateway) *UseCases {
	return &UseCases{ProductOrderGateway: productGateway}
}
