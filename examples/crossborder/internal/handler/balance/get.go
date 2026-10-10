package balance

import (
	"encoding/json"
	"net/http"
	"strings"

	"crossborder/internal/logic/balance"
	"crossborder/internal/svc"
	"crossborder/internal/types"
)

func GetBalanceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetBalanceRequest
		req.SKU = strings.TrimSpace(r.URL.Query().Get("sku"))
		if len(req.SKU) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid sku"})
			return
		}

		l := balance.NewGetBalanceLogic(r.Context(), svcCtx)
		resp, err := l.GetBalance(&req)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
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
