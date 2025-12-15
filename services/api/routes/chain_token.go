package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) ChainTokenApi() {
	r := rs.router
	r.Route("/api/v1/chain-token", func(r chi.Router) {

		r.Post("/create", rs.createChainToken)
		r.Post("/update", rs.updateChainToken)
		r.Get("/info", rs.getChainToken)
		r.Get("/by-chain", rs.listByChain)
		r.Get("/by-token", rs.listByToken)
	})
}

func (rs *Routes) createChainToken(w http.ResponseWriter, r *http.Request) {
	var req service.CreateChainTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.ChainTokenService.CreateChainToken(r.Context(), req)
	if err != nil {
		log.Error("create chain_token failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) updateChainToken(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateChainTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.ChainTokenService.UpdateChainToken(r.Context(), req); err != nil {
		log.Error("update chain_token failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "update success",
	})
}

func (rs *Routes) getChainToken(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.ChainTokenService.GetChainToken(r.Context(), guid)
	if err != nil {
		log.Error("get chain_token failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) listByChain(w http.ResponseWriter, r *http.Request) {
	chainID := r.URL.Query().Get("chain_id")
	if chainID == "" {
		http.Error(w, "chain_id required", http.StatusBadRequest)
		return
	}

	list, err := rs.svc.ChainTokenService.ListByChainID(r.Context(), chainID)
	if err != nil {
		log.Error("list chain_token by chain failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(list)
}

func (rs *Routes) listByToken(w http.ResponseWriter, r *http.Request) {
	tokenID := r.URL.Query().Get("token_id")
	if tokenID == "" {
		http.Error(w, "token_id required", http.StatusBadRequest)
		return
	}

	list, err := rs.svc.ChainTokenService.ListByTokenID(r.Context(), tokenID)
	if err != nil {
		log.Error("list chain_token by token failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(list)
}
