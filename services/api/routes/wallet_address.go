package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) WalletAddressApi() {
	r := rs.router
	r.Route("/api/v1/wallet-address", func(r chi.Router) {

		r.Post("/create", rs.createWalletAddress)
		r.Post("/update", rs.updateWalletAddress)
		r.Get("/info", rs.getWalletAddress)
		r.Get("/by-address", rs.getByAddress)
		r.Get("/by-wallet", rs.listByWalletUUID)
	})
}

func (rs *Routes) createWalletAddress(w http.ResponseWriter, r *http.Request) {
	var req service.CreateWalletAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.WalletAddressService.CreateWalletAddress(r.Context(), req)
	if err != nil {
		log.Error("create wallet_address failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) updateWalletAddress(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateWalletAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.WalletAddressService.UpdateWalletAddress(r.Context(), req); err != nil {
		log.Error("update wallet_address failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "update success",
	})
}

func (rs *Routes) getWalletAddress(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.WalletAddressService.GetWalletAddress(r.Context(), guid)
	if err != nil {
		log.Error("get wallet_address failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) getByAddress(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.WalletAddressService.GetByAddress(r.Context(), address)
	if err != nil {
		log.Error("get wallet_address by address failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) listByWalletUUID(w http.ResponseWriter, r *http.Request) {
	walletUUID := r.URL.Query().Get("wallet_uuid")
	if walletUUID == "" {
		http.Error(w, "wallet_uuid required", http.StatusBadRequest)
		return
	}

	list, err := rs.svc.WalletAddressService.ListByWalletUUID(r.Context(), walletUUID)
	if err != nil {
		log.Error("list wallet_address by wallet failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(list)
}
