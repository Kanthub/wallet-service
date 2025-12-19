package marketapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/roothash-pay/wallet-services/services/market/cache"
)

type Handler struct {
	cache cache.Cache
}

func NewHandler(cache cache.Cache) *Handler {
	return &Handler{cache: cache}
}

func (h *Handler) GetPrice(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	symbol := strings.ToUpper(r.URL.Query().Get("symbol"))
	if symbol == "" {
		http.Error(w, "symbol required", http.StatusBadRequest)
		return
	}

	val, ok, err := h.cache.Get(ctx, "price:"+symbol)
	if err != nil {
		http.Error(w, "cache error", http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "price not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(val)
}
