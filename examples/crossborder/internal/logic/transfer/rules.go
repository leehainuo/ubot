package transfer

import (
	"time"

	"crossborder/internal/model"
)

type transferPlan struct {
	SKU           string
	FromWarehouse string
	ToWarehouse   string
	Quantity      int
	DryRun        bool
}

type transferResult struct {
	Transfer model.Transfer
	From     model.Balance
	To       model.Balance
	Apply    bool
}

func planTransfer(
	plan transferPlan,
	from model.Balance,
	to model.Balance,
	transferID string,
	now time.Time,
) (transferResult, error) {
	if plan.Quantity <= 0 {
		return transferResult{}, ErrInvalidTransfer
	}

	if plan.FromWarehouse == plan.ToWarehouse {
		return transferResult{}, ErrInvalidTransfer
	}

	if from.Available < plan.Quantity {
		return transferResult{}, ErrInsufficientStock
	}

	if plan.DryRun {
		return transferResult{
			Transfer: model.Transfer{
				SKU:           plan.SKU,
				FromWarehouse: plan.FromWarehouse,
				ToWarehouse:   plan.ToWarehouse,
				Quantity:      plan.Quantity,
				Status:        model.TransferValidated,
				CreatedAt:     now,
			},
			From:  from,
			To:    to,
			Apply: false,
		}, nil
	}

	nextFrom := from
	nextTo := to
	nextFrom.Available -= plan.Quantity
	nextTo.Available += plan.Quantity

	return transferResult{
		Transfer: model.Transfer{
			ID:            transferID,
			SKU:           plan.SKU,
			FromWarehouse: plan.FromWarehouse,
			ToWarehouse:   plan.ToWarehouse,
			Quantity:      plan.Quantity,
			Status:        model.TransferCreated,
			CreatedAt:     now,
		},
		From:  nextFrom,
		To:    nextTo,
		Apply: true,
	}, nil
}
