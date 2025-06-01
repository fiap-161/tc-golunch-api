package ports

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/core/dto"
)

type QRCodeProvider interface {
	GenerateQRCode(ctx context.Context, request dto.GenerateQRCodeParams) (string, error)
	VerifyOrder(ctx context.Context, requestUrl string) (any, error)
}
