package inventory

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

type CreateTransferRequest struct {
	SKU            string `json:"sku"`
	FromWarehouse  string `json:"from_warehouse"`
	ToWarehouse    string `json:"to_warehouse"`
	Quantity       int    `json:"quantity"`
	DryRun         bool   `json:"dry_run"`
	IdempotencyKey string `json:"idempotency_key"`
}

type idempotencyRecord struct {
	fingerprint [32]byte
	transfer    Transfer
}

type Service struct {
	mu sync.Mutex

	balances  map[string]Balance
	transfers map[string]Transfer

	idempotency map[string]idempotencyRecord

	nextTransferID int

	now func() time.Time
}

func NewService(now func() time.Time) *Service {
	return &Service{
		balances: map[string]Balance{
			balanceKey("WH_CN_SZ", "SKU_BLACK_01"): {
				WarehouseID: "WH-CN-SZ", SKU: "SKU-BLACK-M-01", Available: 0,
			},
			balanceKey("WH-US-LAX", "SKU-BLACK-M-01"): {
				WarehouseID: "WH-US-LAX", SKU: "SKU-BLACK-M-01", Available: 18,
			},
		},
		transfers:   make(map[string]Transfer),
		idempotency: make(map[string]idempotencyRecord),

		nextTransferID: 1,

		now: now,
	}
}

func (s *Service) SeedBalance(balance Balance) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.balances[balanceKey(balance.WarehouseID, balance.SKU)] = balance
}

func (s *Service) BalanceBySKU(sku string) []Balance {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := make([]Balance, 0, len(s.balances))

	for _, balance := range s.balances {
		if balance.SKU == sku {
			res = append(res, balance)
		}
	}

	return res
}

func (s *Service) CreateTransfer(req CreateTransferRequest) (Transfer, error) {
	// 1. 参数校验
	if req.IdempotencyKey == "" {
		return Transfer{}, ErrIdempotencyKey
	}

	fingerprint := transferFingerprint(req)

	// 2. 幂等校验
	s.mu.Lock()
	defer s.mu.Unlock()

	if cached, ok := s.idempotency[req.IdempotencyKey]; ok {
		if cached.fingerprint != fingerprint {
			return Transfer{}, ErrIdempotencyConflict
		}
	}

	fromKey := balanceKey(req.FromWarehouse, req.SKU)
	toKey := balanceKey(req.ToWarehouse, req.SKU)

	from, ok := s.balances[fromKey]
	if !ok {
		return Transfer{}, ErrNotFound
	}

	to, ok := s.balances[toKey]
	if !ok {
		to = Balance{
			WarehouseID: req.ToWarehouse,
			SKU:         req.SKU,
			Available:   0,
		}
	}

	transferID := s.newTransferIDLocked()
	now := s.now().UTC()

	// 3. 生成 TransferPlan, 并执行计划
	res, err := planTransfer(
		TransferPlan{
			SKU:           req.SKU,
			FromWarehouse: req.FromWarehouse,
			ToWarehouse:   req.ToWarehouse,
			Quantity:      req.Quantity,
			DryRun:        req.DryRun,
		},
		from,
		to,
		transferID,
		now,
	)
	if err != nil {
		return Transfer{}, err
	}

	if !res.Apply {
		return res.Transfer, nil
	}

	// 4. 更新库存和 Transfer 状态
	s.balances[fromKey] = res.From
	s.balances[toKey] = res.To

	s.transfers[transferID] = res.Transfer

	// 5. 缓存幂等记录
	s.idempotency[req.IdempotencyKey] = idempotencyRecord{
		fingerprint: fingerprint,
		transfer:    res.Transfer,
	}

	return res.Transfer, nil
}

func (s *Service) GetTransfer(id string) (Transfer, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	transfer, ok := s.transfers[id]
	if !ok {
		return Transfer{}, false
	}

	return transfer, true
}

func (s *Service) newTransferIDLocked() string {
	id := fmt.Sprintf(
		"TR-%04d",
		s.nextTransferID,
	)

	s.nextTransferID++

	return id
}

func balanceKey(warehouseID, sku string) string {
	return warehouseID + "|" + sku
}

func transferFingerprint(req CreateTransferRequest) [32]byte {
	return sha256.Sum256(fmt.Appendf(nil,
		"%s\x00%s\x00%s\x00%d\x00%t",
		req.SKU,
		req.FromWarehouse,
		req.ToWarehouse,
		req.Quantity,
		req.DryRun,
	))
}
