package order

import (
	"context"

	"crossborder/internal/model"
	"crossborder/internal/svc"
	"crossborder/internal/types"
)

type orderStore interface {
	FindOne(ctx context.Context, id string) (model.Order, error)
}

type GetOrderLogic struct {
	ctx    context.Context
	orders orderStore
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		orders: svcCtx.OrderModel,
	}
}

func (l *GetOrderLogic) GetOrder(req *types.GetOrderRequest) (*types.GetOrderResponse, error) {
	order, err := l.orders.FindOne(l.ctx, req.OrderID)
	if err != nil {
		return nil, err
	}

	return orderToResponse(order), nil
}

func orderToResponse(o model.Order) *types.GetOrderResponse {
	items := make([]types.OrderItem, len(o.Items))
	for i, item := range o.Items {
		items[i] = types.OrderItem{
			SKU:      item.SKU,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	return &types.GetOrderResponse{
		ID:               o.ID,
		Market:           o.Market,
		Currency:         o.Currency,
		Amount:           o.Amount,
		Status:           string(o.Status),
		FulfillmentWH:    o.FulfillmentWH,
		CancellationOpen: o.CancellationOpen,
		Items:            items,
	}
}
