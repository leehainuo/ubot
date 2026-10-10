package order

import (
	"encoding/json"
	"net/http"

	"crossborder/internal/logic/order"
	"crossborder/internal/svc"
	"crossborder/internal/types"
)

func GetOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetOrderRequest
		req.OrderID = r.PathValue("orderID")

		l := order.NewGetOrderLogic(r.Context(), svcCtx)
		resp, err := l.GetOrder(&req)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, resp)
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
