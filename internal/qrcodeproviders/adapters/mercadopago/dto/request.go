package dto

import (
	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/adapters/mercadopago/utils"
	providerdto "github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/core/dto"
)

type RequestGenerateQRCode struct {
	Title             string                      `json:"title"`
	Description       string                      `json:"description"`
	ExternalReference string                      `json:"external_reference"`
	Items             []RequestGenerateQRCodeItem `json:"items"`
	NotificationURL   string                      `json:"notification_url"`
	TotalAmount       float64                     `json:"total_amount"`
}

type RequestGenerateQRCodeItem struct {
	Title       string  `json:"title"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
	UnitMeasure string  `json:"unit_measure"`
	TotalAmount float64 `json:"total_amount"`
	SkuNumber   string  `json:"sku_number,omitempty"`
	Category    string  `json:"category,omitempty"`
	Description string  `json:"description,omitempty"`
}

func (r RequestGenerateQRCode) GetItems() []RequestGenerateQRCodeItem {
	return r.Items
}

func FromParams(params providerdto.GenerateQRCodeParams) RequestGenerateQRCode {
	items, totalAmount := generateItems(params.Items)

	return RequestGenerateQRCode{
		Title:             "Order " + params.OrderID,
		Description:       "Order Description" + params.OrderID,
		ExternalReference: params.OrderID,
		Items:             items,
		NotificationURL:   "https://webhook.site/7694cda3-01d9-4447-a22f-d654af0b9ee2", //TODO adjust this field
		TotalAmount:       utils.FormatDecimal(totalAmount),
	}
}

func generateItems(product []providerdto.Item) ([]RequestGenerateQRCodeItem, float64) {
	items := make([]RequestGenerateQRCodeItem, len(product))
	var totalAmount float64

	for i, item := range product {
		totalAmount += item.Amount
		items[i] = RequestGenerateQRCodeItem{
			Title:       item.Name,
			UnitPrice:   utils.FormatDecimal(item.Price),
			Quantity:    item.Quantity,
			UnitMeasure: "unit",
			TotalAmount: utils.FormatDecimal(item.Amount),
			Description: item.Description,
		}
	}
	return items, totalAmount
}
