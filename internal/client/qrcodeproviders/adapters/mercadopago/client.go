package mercadopago

import (
	"context"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"

	"github.com/fiap-161/tech-challenge-fiap161/internal/client/qrcodeproviders/adapters/mercadopago/dto"
	providerdto "github.com/fiap-161/tech-challenge-fiap161/internal/client/qrcodeproviders/core/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/client/qrcodeproviders/core/ports"
)

type MPClient struct {
	client *resty.Request
}

func New() ports.QRCodeProvider {
	return &MPClient{
		client: getClient(),
	}
}

func getClient() *resty.Request {
	return resty.New().
		R().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + os.Getenv("MERCADO_PAGO_ACCESS_TOKEN"),
		})
}

func (m *MPClient) GenerateQRCode(_ context.Context, params providerdto.GenerateQRCodeParams) (string, error) {
	request := dto.FromParams(params)

	var responseDTO dto.ResponseGenerateQRCode
	_, reqErr := m.client.SetBody(request).SetResult(&responseDTO).Post(viper.GetString("app.mercadopago.url" + "")) //TODO add interpolation rule
	if reqErr != nil {
		return "", reqErr
	}

	return responseDTO.QRData, nil
}
