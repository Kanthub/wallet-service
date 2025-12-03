package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
	"github.com/roothash-pay/wallet-services/services/api/service"
)

// AggregatorRoutes handles swap aggregator related routes
type AggregatorRoutes struct {
	aggregatorService *service.AggregatorService
}

// NewAggregatorRoutes creates a new aggregator routes handler
func NewAggregatorRoutes(aggregatorService *service.AggregatorService) *AggregatorRoutes {
	return &AggregatorRoutes{
		aggregatorService: aggregatorService,
	}
}

// RegisterRoutes registers all aggregator routes to the router
func (h *AggregatorRoutes) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/aggregator", func(r chi.Router) {
		r.Post("/quotes", h.GetQuotesHandler)
		r.Post("/swap/prepare", h.PrepareSwapHandler)
		r.Post("/tx/submitSigned", h.SubmitSignedTxHandler)
		r.Get("/swap/status", h.GetSwapStatusHandler)
	})
}

// GetQuotesHandler handles POST /quotes
func (h *AggregatorRoutes) GetQuotesHandler(w http.ResponseWriter, r *http.Request) {
	var req backend.QuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.aggregatorService.GetQuotes(r.Context(), &req)
	if err != nil {
		log.Error("GetQuotes failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// PrepareSwapHandler handles POST /swap/prepare
func (h *AggregatorRoutes) PrepareSwapHandler(w http.ResponseWriter, r *http.Request) {
	var req backend.PrepareSwapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.aggregatorService.PrepareSwap(r.Context(), &req)
	if err != nil {
		log.Error("PrepareSwap failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// SubmitSignedTxHandler handles POST /tx/submitSigned
func (h *AggregatorRoutes) SubmitSignedTxHandler(w http.ResponseWriter, r *http.Request) {
	var req backend.SubmitSignedTxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.aggregatorService.SubmitSignedTx(r.Context(), &req)
	if err != nil {
		log.Error("SubmitSignedTx failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetSwapStatusHandler handles GET /swap/status
func (h *AggregatorRoutes) GetSwapStatusHandler(w http.ResponseWriter, r *http.Request) {
	swapID := r.URL.Query().Get("swapId")
	if swapID == "" {
		http.Error(w, "swapId is required", http.StatusBadRequest)
		return
	}

	resp, err := h.aggregatorService.GetSwapStatus(r.Context(), swapID)
	if err != nil {
		log.Error("GetSwapStatus failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
