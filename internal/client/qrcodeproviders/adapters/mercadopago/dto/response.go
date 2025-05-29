package dto

type ResponseGenerateQRCode struct {
	InStoreOrderID string `json:"in_store_order_id"`
	QRData         string `json:"qr_data"`
}

type ResponseVerifyOrder struct {
	ID                int64  `json:"id"`
	Status            string `json:"status"`
	ExternalReference string `json:"external_reference"`
	PreferenceID      string `json:"preference_id"`
	OrderStatus       string `json:"order_status"`
}
