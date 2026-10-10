package types

type GetOrderRequest struct {
	OrderID string
}

type GetOrderResponse struct {
	ID               string  `json:"id"`
	Market           string  `json:"market"`
	Currency         string  `json:"currency"`
	Amount           float64 `json:"amount"`
	Status           string  `json:"status"`
	FulfillmentWH    string  `json:"fulfillment_wh"`
	CancellationOpen bool    `json:"cancellation_open"`
	Items            []OrderItem `json:"items"`
}

type OrderItem struct {
	SKU      string  `json:"sku"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
