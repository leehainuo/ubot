package handler

import (
	"net/http"

	"crossborder/internal/handler/balance"
	"crossborder/internal/handler/order"
	"crossborder/internal/handler/transfer"
	"crossborder/internal/svc"
)

func Init(mux *http.ServeMux, svcCtx *svc.ServiceContext) {
	// Ping
	mux.HandleFunc(
		"GET /api/ping",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		})

	// Order
	mux.Handle(
		"GET /api/orders/{orderID}",
		order.GetOrderHandler(svcCtx),
	)

	// Balance
	mux.Handle(
		"GET /api/inventory",
		balance.GetBalanceHandler(svcCtx),
	)

	// Transfer
	mux.Handle(
		"POST /api/transfers",
		transfer.CreateTransferHandler(svcCtx),
	)
}
