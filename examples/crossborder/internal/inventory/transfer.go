package inventory

import "time"

type TransferStatus string

const (
	TransferValidated TransferStatus = "validated"
	TransferCreated   TransferStatus = "created"
)

type Transfer struct {
	ID            string
	SKU           string
	FromWarehouse string
	ToWarehouse   string
	Quantity      int
	Status        TransferStatus
	CreatedAt     time.Time
}

type TransferPlan struct {
	SKU           string
	FromWarehouse string
	ToWarehouse   string
	Quantity      int
	DryRun        bool
}

type TransferResult struct {
	Transfer Transfer

	From Balance
	To   Balance

	Apply bool
}

func planTransfer(
	plan TransferPlan,
	from Balance,
	to Balance,
	transferID string,
	now time.Time,
) (TransferResult, error) {
	// 1. 参数校验
	if plan.Quantity <= 0 {
		return TransferResult{}, ErrInvalidTransfer
	}

	if plan.FromWarehouse == plan.ToWarehouse {
		return TransferResult{}, ErrInvalidTransfer
	}

	if from.Available < plan.Quantity {
		return TransferResult{}, ErrInsufficientStock
	}

	// DryRun 只验证, 不产生实际状态变化
	if plan.DryRun {
		return TransferResult{
			Transfer: Transfer{
				SKU:           plan.SKU,
				FromWarehouse: plan.FromWarehouse,
				ToWarehouse:   plan.ToWarehouse,
				Quantity:      plan.Quantity,
				Status:        TransferValidated,
				CreatedAt:     now,
			},
			From:  from,
			To:    to,
			Apply: false,
		}, nil
	}

	// 2. 执行 Transfer
	nextFrom := from
	nextTo := to
	nextFrom.Available -= plan.Quantity
	nextTo.Available += plan.Quantity

	return TransferResult{
		Transfer: Transfer{
			ID:            transferID,
			SKU:           plan.SKU,
			FromWarehouse: plan.FromWarehouse,
			ToWarehouse:   plan.ToWarehouse,
			Quantity:      plan.Quantity,
			Status:        TransferCreated,
			CreatedAt:     now,
		},
		From:     nextFrom,
		To:       nextTo,
		Apply: true,
	}, nil
}
