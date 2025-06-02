package ports

import (
	"context"
	dto2 "github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/adapters/mercadopago/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/core/dto"
)

type QRCodeProvider interface {
	GenerateQRCode(ctx context.Context, request dto.GenerateQRCodeParams) (string, error)
	CheckPayment(ctx context.Context, requestUrl string) (dto2.ResponseVerifyOrder, error)
}
