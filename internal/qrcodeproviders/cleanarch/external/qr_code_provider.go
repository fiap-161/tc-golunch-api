package external

import (
	"context"

	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/cleanarch/entity"
)

type QRCodeProvider interface {
	GenerateQRCode(ctx context.Context, request entity.GenerateQRCodeParams) (string, error)
	CheckPayment(ctx context.Context, requestUrl string) (entity.ResponseVerifyOrder, error)
}
