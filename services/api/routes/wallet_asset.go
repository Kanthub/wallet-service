package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) WalletAssetApi() {
	r := rs.router
	r.Route("/api/v1/wallet-asset", func(r chi.Router) {

		r.Post("/create", rs.createWalletAsset)
		r.Post("/update", rs.updateWalletAsset)
		r.Get("/info", rs.getWalletAsset)
		r.Get("/by-token-chain", rs.getByTokenChain)
	})
}

func (rs *Routes) createWalletAsset(w http.ResponseWriter, r *http.Request) {
	var req service.CreateWalletAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.WalletAssetService.CreateWalletAsset(r.Context(), req)
	if err != nil {
		log.Error("create wallet_asset failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) updateWalletAsset(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateWalletAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.WalletAssetService.UpdateWalletAsset(r.Context(), req); err != nil {
		log.Error("update wallet_asset failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "update success",
	})
}

func (rs *Routes) getWalletAsset(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.WalletAssetService.GetWalletAsset(r.Context(), guid)
	if err != nil {
		log.Error("get wallet_asset failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) getByTokenChain(w http.ResponseWriter, r *http.Request) {
	tokenID := r.URL.Query().Get("token_id")
	chainID := r.URL.Query().Get("chain_id")

	if tokenID == "" || chainID == "" {
		http.Error(w, "token_id and chain_id required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.WalletAssetService.GetByTokenChain(r.Context(), tokenID, chainID)
	if err != nil {
		log.Error("get wallet_asset by token+chain failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}
