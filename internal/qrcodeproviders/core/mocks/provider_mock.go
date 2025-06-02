package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/core/dto"
)

type MockQRCodeProvider struct {
	mock.Mock
}

func (m *MockQRCodeProvider) GenerateQRCode(ctx context.Context, request dto.GenerateQRCodeParams) (string, error) {
	args := m.Called(ctx, request)
	return args.String(0), args.Error(1)
}

func (m *MockQRCodeProvider) CheckPayment(ctx context.Context, requestUrl string) (any, error) {
	args := m.Called(ctx, requestUrl)
	return args.Get(0), args.Error(1)
}
