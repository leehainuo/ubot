package svc

import "crossborder/internal/model/memory"

type ServiceContext struct {
	OrderModel    *memory.OrderModel
	BalanceModel  *memory.BalanceModel
	TransferModel *memory.TransferModel
}

func NewServiceContext() *ServiceContext {
	return &ServiceContext{
		OrderModel:    memory.NewOrderModel(),
		BalanceModel:  memory.NewBalanceModel(),
		TransferModel: memory.NewTransferModel(),
	}
}
