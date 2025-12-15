package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"
)

// 注册路由
func (rs *Routes) FiatCurrencyRateApi() {
	r := rs.router

	r.Route("/api/v1/fiat-currency-rate", func(r chi.Router) {
		r.Post("/set", rs.setFiatCurrencyRate)
		r.Get("/get", rs.getFiatCurrencyRate)
	})
}

/*
POST /api/v1/fiat-currency-rate/set

	{
		"key": "USD_CNY",
		"value": "7.25"
	}
*/
type setFiatCurrencyRateRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (rs *Routes) setFiatCurrencyRate(w http.ResponseWriter, r *http.Request) {
	var req setFiatCurrencyRateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.FiatCurrencyRateService.SetRate(
		r.Context(),
		req.Key,
		req.Value,
	); err != nil {
		log.Error("set fiat currency rate failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "set rate success",
	})
}

/*
GET /api/v1/fiat-currency-rate/get?key=USD_CNY
*/
func (rs *Routes) getFiatCurrencyRate(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.FiatCurrencyRateService.GetRate(
		r.Context(),
		key,
	)
	if err != nil {
		log.Error("get fiat currency rate failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}
