package model

type OrderStatus string

const (
	StatusAwaitingShipment OrderStatus = "awaiting_shipment"
	StatusCancelled        OrderStatus = "cancelled"
)

type OrderItem struct {
	SKU      string
	Quantity int
	Price    float64
}

type Order struct {
	ID               string
	Market           string
	Currency         string
	Amount           float64
	Status           OrderStatus
	FulfillmentWH    string
	CancellationOpen bool
	Items            []OrderItem
}
