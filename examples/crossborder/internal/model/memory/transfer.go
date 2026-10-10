package memory

import (
	"context"
	"fmt"
	"sync"

	"crossborder/internal/model"
)

type TransferModel struct {
	mu sync.Mutex

	transfers   map[string]model.Transfer
	idempotency map[string]model.IdempotencyRecord

	nextTransferID int
}

func NewTransferModel() *TransferModel {
	return &TransferModel{
		transfers:      make(map[string]model.Transfer),
		idempotency:    make(map[string]model.IdempotencyRecord),
		nextTransferID: 1,
	}
}

func (m *TransferModel) NextID(_ context.Context) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := fmt.Sprintf("TR-%04d", m.nextTransferID)
	m.nextTransferID++
	return id, nil
}

func (m *TransferModel) Create(_ context.Context, t model.Transfer) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.transfers[t.ID] = t
	return nil
}

func (m *TransferModel) FindOne(_ context.Context, id string) (model.Transfer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.transfers[id]
	if !ok {
		return model.Transfer{}, model.ErrNotFound
	}
	return t, nil
}

func (m *TransferModel) GetIdempotency(_ context.Context, key string) (model.IdempotencyRecord, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	rec, ok := m.idempotency[key]
	if !ok {
		return model.IdempotencyRecord{}, model.ErrNotFound
	}
	return rec, nil
}

func (m *TransferModel) SetIdempotency(_ context.Context, key string, rec model.IdempotencyRecord) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.idempotency[key] = rec
	return nil
}
