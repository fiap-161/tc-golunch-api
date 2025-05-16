package mercadopago

import (
	"context"
	"github.com/fiap-161/tech-challenge-fiap161/internal/client/mercadopago/dto"
	"github.com/go-resty/resty/v2"
	"os"
)

type MercadoPago struct {
	defaultClient *resty.Request
}

type QRCodeProvider interface {
	GenerateQRCode(ctx context.Context, request dto.RequestGenerateQRCode) (string, error) // QRCode, err
}

func New() QRCodeProvider {
	return &MercadoPago{
		defaultClient: getClient(),
	}
}

func getClient() *resty.Request {
	return resty.New().
		R().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + os.Getenv("MERCADOPAGO_ACCESS_TOKEN"),
		})
}

func (m *MercadoPago) GenerateQRCode(_ context.Context, request dto.RequestGenerateQRCode) (string, error) {
	
}
