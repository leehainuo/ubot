package httpapi

import (
	"crossborder/internal/inventory"
	"crossborder/internal/order"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type OrderService interface {
	GetOrder(id string) (order.Order, bool)
}

type InventoryService interface {
	BalanceBySKU(sku string) []inventory.Balance
	CreateTransfer(req inventory.CreateTransferRequest) (inventory.Transfer, error)
}

func New(o OrderService, i InventoryService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// 注册业务相关接口

	// 1. 查订单
	mux.HandleFunc("GET /api/orders/{orderID}", func(w http.ResponseWriter, r *http.Request) {
		order, ok := o.GetOrder(r.PathValue("orderID"))
		if !ok {
			// 返回错误
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "order not found"})
			return
		}
		writeJSON(w, http.StatusOK, order)
	})

	// 2. 查库存
	mux.HandleFunc("GET /api/inventory", func(w http.ResponseWriter, r *http.Request) {
		sku := strings.TrimSpace(r.URL.Query().Get("sku"))
		if len(sku) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid sku"})
			return
		}
		data := i.BalanceBySKU(sku)
		writeJSON(w, http.StatusOK, data)
	})

	// 3. 创建仓库调拨单
	mux.HandleFunc("POST /api/transfers", func(w http.ResponseWriter, r *http.Request) {
		// 解析请求参数
		var req inventory.CreateTransferRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err})
			return
		}
		transfer, err := i.CreateTransfer(req)
		if err != nil {
			if errors.Is(err, inventory.ErrIdempotencyKey) || errors.Is(err, inventory.ErrIdempotencyConflict) {
				writeJSON(w, http.StatusBadRequest, map[string]any{"error": err})
				return
			}
			writeJSON(w, http.StatusConflict, map[string]any{"error": err})
			return
		}
		writeJSON(w, http.StatusOK, transfer)
	})

	return mux
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
