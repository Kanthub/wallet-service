package routes

import (
	"encoding/json"
	"net/http"
	"strings"

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

		r.Get("/quote", rs.getMarketQuote)
		r.Get("/quotes", rs.getMarketQuotes)

	})
}

// getMarketQuote godoc
// @Summary Get real-time market price
// @Description Get latest market price from cache (info-level)
// @Tags Market
// @Accept json
// @Produce json
// @Param symbol query string true "Asset symbol, e.g. BTC"
// @Success 200 {object} model.Quote
// @Failure 400 {string} string "symbol required"
// @Failure 404 {string} string "price not found"
// @Router /api/v1/market-price/quote [get]
func (rs *Routes) getMarketQuote(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "symbol required", http.StatusBadRequest)
		return
	}

	quote, err := rs.svc.MarketPriceService.GetPrice(
		r.Context(),
		symbol,
	)
	if err != nil {
		log.Warn("get market quote failed", "symbol", symbol, "err", err)
		http.Error(w, "price not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(quote)
}

// getMarketQuotes godoc
// @Summary Get batch market prices
// @Description Get latest market prices for multiple symbols from cache
// @Tags Market
// @Accept json
// @Produce json
// @Param symbols query string true "Comma-separated symbols, e.g. BTC,ETH,SOL"
// @Success 200 {object} map[string]model.Quote
// @Failure 400 {string} string "symbols required"
// @Router /api/v1/market-price/quotes [get]
func (rs *Routes) getMarketQuotes(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.Query().Get("symbols")
	if raw == "" {
		http.Error(w, "symbols required", http.StatusBadRequest)
		return
	}

	symbols := strings.Split(raw, ",")
	result := make(map[string]interface{})

	for _, s := range symbols {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		q, err := rs.svc.MarketPriceService.GetPrice(r.Context(), s)
		if err != nil {
			continue
		}
		result[strings.ToUpper(s)] = q
	}

	json.NewEncoder(w).Encode(result)
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
