package memory

import (
	"context"
	"sync"

	"crossborder/internal/model"
)

type OrderModel struct {
	mu     sync.RWMutex
	orders map[string]model.Order
}

func NewOrderModel() *OrderModel {
	return &OrderModel{
		orders: map[string]model.Order{
			"TTS-20260801-1001": {
				ID:               "TTS-20260801-1001",
				Market:           "US",
				Currency:         "USD",
				Amount:           129.99,
				Status:           model.StatusAwaitingShipment,
				FulfillmentWH:    "WH-CN-SZ",
				CancellationOpen: true,
				Items:            []model.OrderItem{{SKU: "SKU-BLACK-M-01", Quantity: 1, Price: 129.99}},
			},
		},
	}
}

func (m *OrderModel) Seed(o model.Order) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.orders[o.ID] = o
}

func (m *OrderModel) FindOne(_ context.Context, id string) (model.Order, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	o, ok := m.orders[id]
	if !ok {
		return model.Order{}, model.ErrNotFound
	}
	return o, nil
}
