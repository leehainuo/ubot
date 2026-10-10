package types

type GetBalanceRequest struct {
	SKU string `json:"sku"`
}

type GetBalanceResponse struct {
	WarehouseID string `json:"warehouse_id"`
	SKU         string `json:"sku"`
	Available   int    `json:"available"`
}
