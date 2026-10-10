package transfer

import (
	"encoding/json"
	"errors"
	"net/http"

	"crossborder/internal/logic/transfer"
	"crossborder/internal/svc"
	"crossborder/internal/types"
)

func CreateTransferHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateTransferRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}

		l := transfer.NewCreateTransferLogic(r.Context(), svcCtx)
		resp, err := l.CreateTransfer(&req)
		if err != nil {
			if errors.Is(err, transfer.ErrIdempotencyKey) ||
				errors.Is(err, transfer.ErrIdempotencyConflict) ||
				errors.Is(err, transfer.ErrInvalidTransfer) {
				writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusConflict, map[string]any{"error": err.Error()})
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
