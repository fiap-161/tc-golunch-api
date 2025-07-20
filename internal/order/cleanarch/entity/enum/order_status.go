package enum

type OrderStatus string

const (
	OrderStatusAwaitingPayment OrderStatus = "awaiting_payment"
	OrderStatusReceived        OrderStatus = "received"
	OrderStatusInPreparation   OrderStatus = "in_preparation"
	OrderStatusReady           OrderStatus = "ready"
	OrderStatusCompleted       OrderStatus = "completed"
)

var OrderPanelStatus = []string{
	OrderStatusReceived.String(),
	OrderStatusInPreparation.String(),
	OrderStatusReady.String(),
}

func (o OrderStatus) String() string {
	return string(o)
}
