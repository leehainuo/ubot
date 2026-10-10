package memory

import (
	"context"
	"sync"

	"crossborder/internal/model"
)

type BalanceModel struct {
	mu       sync.RWMutex
	balances map[string]model.Balance
}

func NewBalanceModel() *BalanceModel {
	return &BalanceModel{
		balances: map[string]model.Balance{
			balanceKey("WH-CN-SZ", "SKU-BLACK-M-01"): {
				WarehouseID: "WH-CN-SZ",
				SKU:         "SKU-BLACK-M-01",
				Available:   0,
			},
			balanceKey("WH-US-LAX", "SKU-BLACK-M-01"): {
				WarehouseID: "WH-US-LAX",
				SKU:         "SKU-BLACK-M-01",
				Available:   18,
			},
		},
	}
}

func (m *BalanceModel) Seed(b model.Balance) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.balances[balanceKey(b.WarehouseID, b.SKU)] = b
}

func (m *BalanceModel) FindOne(_ context.Context, sku string, warehouseID string) (model.Balance, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	b, ok := m.balances[balanceKey(warehouseID, sku)]
	if !ok {
		return model.Balance{}, model.ErrNotFound
	}
	return b, nil
}

func (m *BalanceModel) FindBySKU(_ context.Context, sku string) ([]model.Balance, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]model.Balance, 0)
	for _, b := range m.balances {
		if b.SKU == sku {
			res = append(res, b)
		}
	}
	return res, nil
}

func (m *BalanceModel) Update(_ context.Context, b model.Balance) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.balances[balanceKey(b.WarehouseID, b.SKU)] = b
	return nil
}

func balanceKey(warehouseID, sku string) string {
	return warehouseID + "|" + sku
}
