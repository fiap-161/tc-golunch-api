package external

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/cleanarch/dtos"
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/cleanarch/entities"
)

type QRCodeProvider interface {
	GenerateQRCode(ctx context.Context, request entities.GenerateQRCodeParams) (string, error)
	CheckPayment(ctx context.Context, requestUrl string) (dtos.ResponseVerifyOrderDTO, error)
}
