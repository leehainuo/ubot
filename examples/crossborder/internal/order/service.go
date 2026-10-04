package order

import "sync"

type Service struct {
	mu     sync.Mutex
	orders map[string]Order
}

func NewService() *Service {
	return &Service{
		orders: map[string]Order{
			"TTS-20260801-1001": {
				ID:               "TTS-20260801-1001",
				Market:           "US",
				Currency:         "USD",
				Amount:           129.99,
				Status:           StatusAwaitingShipment,
				FulfillmentWH:    "WH-CN-SZ",
				CancellationOpen: true,
				Items:            []Item{{SKU: "SKU-BLACK-M-01", Quantity: 1, Price: 129.99}},
			},
		},
	}
}

func (s *Service) Seed(o Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[o.ID] = o
}

func (s *Service) GetOrder(id string) (Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	o, ok := s.orders[id]
	if !ok {
		return Order{}, false
	}

	return o, true
}
