package dto

type GenerateQRCodeParams struct {
	OrderID string
	Items   []Product
}

// TODO remove it after integrate with payment
type Product struct {
	ID            string
	Name          string
	Category      string
	Price         float64
	Description   string
	ImageURL      string
	PreparingTime int
	Quantity      int     // TODO verify
	TotalPrice    float64 // TODO verify
}
