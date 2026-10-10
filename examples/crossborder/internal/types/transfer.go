package types

type CreateTransferRequest struct {
	SKU            string `json:"sku"`
	FromWarehouse  string `json:"from_warehouse"`
	ToWarehouse    string `json:"to_warehouse"`
	Quantity       int    `json:"quantity"`
	DryRun         bool   `json:"dry_run"`
	IdempotencyKey string `json:"idempotency_key"`
}

type CreateTransferResponse struct {
	ID            string `json:"id"`
	SKU           string `json:"sku"`
	FromWarehouse string `json:"from_warehouse"`
	ToWarehouse   string `json:"to_warehouse"`
	Quantity      int    `json:"quantity"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}
