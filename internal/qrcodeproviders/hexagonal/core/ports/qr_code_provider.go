package ports

import (
	"context"

	dto2 "github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/hexagonal/adapters/mercadopago/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/hexagonal/core/dto"
)

type QRCodeProvider interface {
	GenerateQRCode(ctx context.Context, request dto.GenerateQRCodeParams) (string, error)
	CheckPayment(ctx context.Context, requestUrl string) (dto2.ResponseVerifyOrder, error)
}
