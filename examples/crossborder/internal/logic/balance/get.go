package balance

import (
	"context"

	"crossborder/internal/model"
	"crossborder/internal/svc"
	"crossborder/internal/types"
)

type balanceReader interface {
	FindBySKU(ctx context.Context, sku string) ([]model.Balance, error)
}

type GetBalanceLogic struct {
	ctx      context.Context
	balances balanceReader
}

func NewGetBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBalanceLogic {
	return &GetBalanceLogic{
		ctx:      ctx,
		balances: svcCtx.BalanceModel,
	}
}

func (l *GetBalanceLogic) GetBalance(req *types.GetBalanceRequest) ([]types.GetBalanceResponse, error) {
	balances, err := l.balances.FindBySKU(l.ctx, req.SKU)
	if err != nil {
		return nil, err
	}

	result := make([]types.GetBalanceResponse, len(balances))
	for i, b := range balances {
		result[i] = types.GetBalanceResponse{
			WarehouseID: b.WarehouseID,
			SKU:         b.SKU,
			Available:   b.Available,
		}
	}

	return result, nil
}
