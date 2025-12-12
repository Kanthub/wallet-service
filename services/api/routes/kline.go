package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) KlineApi() {
	r := rs.router

	r.Route("/api/v1/kline", func(r chi.Router) {
		r.Post("/set", rs.setKline)
		r.Get("/list", rs.getKlines)
	})
}

func (rs *Routes) setKline(w http.ResponseWriter, r *http.Request) {
	var req service.SetKlineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.KlineService.SetKline(
		r.Context(),
		req,
	); err != nil {
		log.Error("set kline failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "set kline success",
	})
}

func (rs *Routes) getKlines(w http.ResponseWriter, r *http.Request) {
	tokenID := r.URL.Query().Get("token_id")
	interval := r.URL.Query().Get("interval")

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	limitStr := r.URL.Query().Get("limit")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			http.Error(w, "invalid start time", http.StatusBadRequest)
			return
		}
	}
	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			http.Error(w, "invalid end time", http.StatusBadRequest)
			return
		}
	}

	limit := 500
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	list, err := rs.svc.KlineService.GetKlines(
		r.Context(),
		tokenID,
		interval,
		start,
		end,
		limit,
	)
	if err != nil {
		log.Error("get klines failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(list)
}
