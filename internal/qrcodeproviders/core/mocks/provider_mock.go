package mocks

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/core/dto"
	"github.com/stretchr/testify/mock"
)

type MockQRCodeProvider struct {
	mock.Mock
}

func (m *MockQRCodeProvider) GenerateQRCode(ctx context.Context, request dto.GenerateQRCodeParams) (string, error) {
	args := m.Called(ctx, request)
	return args.String(0), args.Error(1)
}
