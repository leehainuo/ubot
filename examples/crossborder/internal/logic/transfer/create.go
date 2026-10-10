package transfer

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"crossborder/internal/model"
	"crossborder/internal/svc"
	"crossborder/internal/types"
)

type balanceStore interface {
	FindOne(ctx context.Context, sku, warehouseID string) (model.Balance, error)
	Update(ctx context.Context, b model.Balance) error
}

type transferStore interface {
	NextID(ctx context.Context) (string, error)
	Create(ctx context.Context, t model.Transfer) error
	GetIdempotency(ctx context.Context, key string) (model.IdempotencyRecord, error)
	SetIdempotency(ctx context.Context, key string, rec model.IdempotencyRecord) error
}

type CreateTransferLogic struct {
	ctx       context.Context
	balances  balanceStore
	transfers transferStore
	now       func() time.Time
}

func NewCreateTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTransferLogic {
	return &CreateTransferLogic{
		ctx:       ctx,
		balances:  svcCtx.BalanceModel,
		transfers: svcCtx.TransferModel,
		now:       time.Now,
	}
}

func (l *CreateTransferLogic) CreateTransfer(req *types.CreateTransferRequest) (*types.CreateTransferResponse, error) {
	// 1. 参数校验
	if req.IdempotencyKey == "" {
		return nil, ErrIdempotencyKey
	}

	fingerprint := transferFingerprint(req)

	// 2. 幂等校验
	rec, err := l.transfers.GetIdempotency(l.ctx, req.IdempotencyKey)
	switch {
	case err == nil:
		if rec.Fingerprint != fingerprint {
			return nil, ErrIdempotencyConflict
		}
		return transferToResponse(rec.Transfer), nil
	case !errors.Is(err, model.ErrNotFound):
		return nil, err
	}

	// 3. 查询库存
	from, err := l.balances.FindOne(l.ctx, req.SKU, req.FromWarehouse)
	if err != nil {
		return nil, err
	}

	to, err := l.balances.FindOne(l.ctx, req.SKU, req.ToWarehouse)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			return nil, err
		}
		to = model.Balance{
			WarehouseID: req.ToWarehouse,
			SKU:         req.SKU,
			Available:   0,
		}
	}

	transferID, err := l.transfers.NextID(l.ctx)
	if err != nil {
		return nil, err
	}
	now := l.now().UTC()

	// 4. 调用纯函数计算调拨结果
	res, err := planTransfer(
		transferPlan{
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
		return nil, err
	}

	if !res.Apply {
		return transferToResponse(res.Transfer), nil
	}

	// 5. 持久化
	if err := l.balances.Update(l.ctx, res.From); err != nil {
		return nil, err
	}
	if err := l.balances.Update(l.ctx, res.To); err != nil {
		return nil, err
	}
	if err := l.transfers.Create(l.ctx, res.Transfer); err != nil {
		return nil, err
	}
	if err := l.transfers.SetIdempotency(l.ctx, req.IdempotencyKey, model.IdempotencyRecord{
		Fingerprint: fingerprint,
		Transfer:    res.Transfer,
	}); err != nil {
		return nil, err
	}

	return transferToResponse(res.Transfer), nil
}

func transferFingerprint(req *types.CreateTransferRequest) [32]byte {
	return sha256.Sum256(fmt.Appendf(nil,
		"%s\x00%s\x00%s\x00%d\x00%t",
		req.SKU,
		req.FromWarehouse,
		req.ToWarehouse,
		req.Quantity,
		req.DryRun,
	))
}

func transferToResponse(t model.Transfer) *types.CreateTransferResponse {
	return &types.CreateTransferResponse{
		ID:            t.ID,
		SKU:           t.SKU,
		FromWarehouse: t.FromWarehouse,
		ToWarehouse:   t.ToWarehouse,
		Quantity:      t.Quantity,
		Status:        string(t.Status),
		CreatedAt:     t.CreatedAt.Format(time.RFC3339),
	}
}
