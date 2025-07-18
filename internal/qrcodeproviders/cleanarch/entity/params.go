package entity

type GenerateQRCodeParams struct {
	OrderID string
	Items   []Item
}

type Item struct {
	ID          string
	Name        string
	Price       float64
	Description string
	Quantity    int
	Amount      float64
}

type VerifyOrderResponse struct {
	ExternalReference string `json:"external_reference"`
	OrderStatus       string `json:"order_status"`
}
