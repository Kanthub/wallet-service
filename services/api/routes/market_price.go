package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) MarketPriceApi() {
	r := rs.router

	r.Route("/api/v1/market-price", func(r chi.Router) {
		r.Post("/set", rs.setMarketPrice)
		r.Get("/by-token", rs.getMarketPriceByToken)
		r.Get("/info", rs.getMarketPriceByGuid)
	})
}

func (rs *Routes) setMarketPrice(w http.ResponseWriter, r *http.Request) {
	var req service.SetMarketPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.MarketPriceService.SetMarketPrice(
		r.Context(),
		req,
	); err != nil {
		log.Error("set market price failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "set market price success",
	})
}

func (rs *Routes) getMarketPriceByToken(w http.ResponseWriter, r *http.Request) {
	tokenID := r.URL.Query().Get("token_id")
	if tokenID == "" {
		http.Error(w, "token_id required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.MarketPriceService.GetByTokenID(
		r.Context(),
		tokenID,
	)
	if err != nil {
		log.Error("get market price by token failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) getMarketPriceByGuid(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.MarketPriceService.GetByGuid(
		r.Context(),
		guid,
	)
	if err != nil {
		log.Error("get market price failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}
