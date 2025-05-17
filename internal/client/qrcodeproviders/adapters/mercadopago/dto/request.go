package dto

import (
	"time"

	providerdto "github.com/fiap-161/tech-challenge-fiap161/internal/client/qrcodeproviders/core/dto"
)

type RequestGenerateQRCode struct {
	Title             string                      `json:"title"`
	Description       string                      `json:"description"`
	ExpirationDate    string                      `json:"expiration_date"`
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

func FromParams(params providerdto.GenerateQRCodeParams) RequestGenerateQRCode {
	items, totalAmount := fromProducts(params.Items)

	return RequestGenerateQRCode{
		Title:             "Order " + params.OrderID,
		Description:       "Order " + params.OrderID, // TODO verify
		ExpirationDate:    time.Now().Add(10 * time.Minute).Format(time.RFC3339),
		ExternalReference: params.OrderID,
		Items:             items,
		NotificationURL:   "https://example.com/notifications", //TODO edit
		TotalAmount:       totalAmount,
	}
}

func fromProducts(product []providerdto.Product) ([]RequestGenerateQRCodeItem, float64) {
	items := make([]RequestGenerateQRCodeItem, len(product))
	var totalAmount float64

	for _, item := range product {
		totalAmount += item.TotalPrice
		items = append(items, RequestGenerateQRCodeItem{
			Title:       item.Name,
			UnitPrice:   item.Price,
			Quantity:    item.Quantity,
			UnitMeasure: "unit", // TODO verify
			TotalAmount: item.TotalPrice,
			SkuNumber:   "", // TODO verify
			Category:    "", // TODO verify
			Description: item.Description,
		})
	}
	return items, totalAmount
}

func (r RequestGenerateQRCode) GetItems() []RequestGenerateQRCodeItem {
	return r.Items
}
