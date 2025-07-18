package adapters

import (
	"os"

	"github.com/fiap-161/tech-challenge-fiap161/internal/qrcodeproviders/cleanarch/entity"
)

func FromParams(params entity.GenerateQRCodeParams) entity.RequestGenerateQRCode {
	items, totalAmount := generateItems(params.Items)

	return entity.RequestGenerateQRCode{
		Title:             "Order " + params.OrderID,
		Description:       "Order Description" + params.OrderID,
		ExternalReference: params.OrderID,
		Items:             items,
		NotificationURL:   os.Getenv("WEBHOOK_URL"), //TODO adjust this field
		TotalAmount:       FormatDecimal(totalAmount),
	}
}

func generateItems(product []entity.Item) ([]entity.RequestGenerateQRCodeItem, float64) {
	items := make([]entity.RequestGenerateQRCodeItem, len(product))
	var totalAmount float64

	for i, item := range product {
		totalAmount += item.Amount
		items[i] = entity.RequestGenerateQRCodeItem{
			Title:       item.Name,
			UnitPrice:   FormatDecimal(item.Price),
			Quantity:    item.Quantity,
			UnitMeasure: "unit",
			TotalAmount: FormatDecimal(item.Amount),
			Description: item.Description,
		}
	}
	return items, totalAmount
}
