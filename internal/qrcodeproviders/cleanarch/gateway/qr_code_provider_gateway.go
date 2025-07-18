package gateway

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/cleanarch/adapters"
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/cleanarch/entity"
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/cleanarch/external"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"

	"github.com/fiap-161/tech-challenge-fiap161/internal/shared"
)

var (
	SellerUserID  = os.Getenv("MERCADO_PAGO_SELLER_APP_USER_ID")
	ExternalPosID = os.Getenv("MERCADO_PAGO_EXTERNAL_POS_ID")
)

type MercadoPagoClient struct {
	client *resty.Client
}

func New() external.QRCodeProvider {
	return &MercadoPagoClient{
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

func (m *MercadoPagoClient) GenerateQRCode(_ context.Context, params entity.GenerateQRCodeParams) (string, error) {
	request := adapters.FromParams(params)

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

	var responseDTO entity.ResponseGenerateQRCode
	res, reqErr := m.client.R().
		SetBody(request).
		SetResult(&responseDTO).
		Post(resolvedPath)

	if res != nil && res.IsError() {
		fmt.Println(res.Error())
		return "", errors.New("error in request, endpoint called: " + res.Request.URL)
	}
	if reqErr != nil {
		return "", reqErr
	}

	return responseDTO.QRData, nil
}

func (m *MercadoPagoClient) CheckPayment(_ context.Context, requestUrl string) (entity.ResponseVerifyOrder, error) {
	var responseDTO entity.ResponseVerifyOrder
	res, reqErr := m.client.R().
		SetResult(&responseDTO).
		Get(requestUrl)
	if res != nil && res.IsError() {
		return entity.ResponseVerifyOrder{}, errors.New("error in request, endpoint called: " + res.Request.URL + "\n")
	}
	if reqErr != nil {
		return entity.ResponseVerifyOrder{}, reqErr
	}

	return responseDTO, nil
}
