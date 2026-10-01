package order

type Status string

const (
	StatusAwaitingShipment Status = "awaiting_shipment"
	StatusCancelled        Status = "cancelled"
)

type Item struct {
	SKU      string
	Quantity int
	Price    float64
}

type Order struct {
	ID               string
	Market           string
	Currency         string
	Amount           float64
	Status           Status
	FulfillmentWH    string
	CancellationOpen bool
	Items            []Item
}
