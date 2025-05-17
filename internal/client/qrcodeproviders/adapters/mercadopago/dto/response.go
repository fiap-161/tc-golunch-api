package dto

type ResponseGenerateQRCode struct {
	InStoreOrderID string `json:"in_store_order_id"`
	QRData         string `json:"qr_data"`
}
