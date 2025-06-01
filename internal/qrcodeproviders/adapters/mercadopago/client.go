package mercadopago

import (
	"context"
	"errors"
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/adapters/mercadopago/dto"
	providerdto "github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/core/dto"
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/core/ports"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"

	"github.com/fiap-161/tech-challenge-fiap161/internal/shared"
)

var (
	SellerUserID  = os.Getenv("MERCADO_PAGO_SELLER_APP_USER_ID")
	ExternalPosID = os.Getenv("MERCADO_PAGO_EXTERNAL_POS_ID")
)

type MPClient struct {
	client *resty.Client
}

func New() ports.QRCodeProvider {
	return &MPClient{
		client: getClient(),
	}
}

func getClient() *resty.Client {
	return resty.New().
		SetBaseURL(viper.GetString(shared.MercadoPagoHost)).
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + os.Getenv("MERCADO_PAGO_ACCESS_TOKEN"),
		})
}

func (m *MPClient) GenerateQRCode(_ context.Context, params providerdto.GenerateQRCodeParams) (string, error) {
	request := dto.FromParams(params)

	pathParams := []shared.BuildPathParam{
		{
			Key:   "user_id",
			Value: SellerUserID,
		},
		{
			Key:   "external_pos_id",
			Value: ExternalPosID,
		},
	}
	resolvedPath, err := shared.BuildPath(viper.GetString(shared.MercadoPagoQRCodePath), pathParams)
	if err != nil {
		return "", err
	}

	var responseDTO dto.ResponseGenerateQRCode
	res, reqErr := m.client.R().
		SetBody(request).
		SetResult(&responseDTO).
		Post(resolvedPath)
	if res != nil && res.IsError() {
		return "", errors.New("error in request, endpoint called: " + res.Request.URL)
	}
	if reqErr != nil {
		return "", reqErr
	}

	return responseDTO.QRData, nil
}

func (m *MPClient) CheckPayment(_ context.Context, requestUrl string) (any, error) {
	var responseDTO dto.ResponseVerifyOrder
	res, reqErr := m.client.R().
		SetResult(&responseDTO).
		Get(requestUrl)
	if res != nil && res.IsError() {
		return nil, errors.New("error in request, endpoint called: " + res.Request.URL + "\n")
	}
	if reqErr != nil {
		return nil, reqErr
	}

	return responseDTO, nil
}
