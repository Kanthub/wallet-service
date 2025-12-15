package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) WalletApi() {
	r := rs.router
	r.Route("/api/v1/wallet", func(r chi.Router) {

		r.Post("/create", rs.createWallet)
		r.Post("/update", rs.updateWallet)
		r.Get("/info", rs.getWallet)
		r.Get("/by-uuid", rs.getWalletByUUID)
		r.Get("/list", rs.listWallets)
	})
}

func (rs *Routes) createWallet(w http.ResponseWriter, r *http.Request) {
	var req service.CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	wallet, err := rs.svc.WalletService.CreateWallet(r.Context(), req)
	if err != nil {
		log.Error("create wallet failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(wallet)
}

func (rs *Routes) updateWallet(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.WalletService.UpdateWallet(r.Context(), req); err != nil {
		log.Error("update wallet failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "update success",
	})
}

func (rs *Routes) getWallet(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid required", http.StatusBadRequest)
		return
	}

	wallet, err := rs.svc.WalletService.GetWallet(r.Context(), guid)
	if err != nil {
		log.Error("get wallet failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(wallet)
}

func (rs *Routes) getWalletByUUID(w http.ResponseWriter, r *http.Request) {
	walletUUID := r.URL.Query().Get("wallet_uuid")
	if walletUUID == "" {
		http.Error(w, "wallet_uuid required", http.StatusBadRequest)
		return
	}

	wallet, err := rs.svc.WalletService.GetByWalletUUID(r.Context(), walletUUID)
	if err != nil {
		log.Error("get wallet by uuid failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(wallet)
}

func (rs *Routes) listWallets(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	filters := map[string]interface{}{
		"wallet_name": r.URL.Query().Get("wallet_name"),
		"device_uuid": r.URL.Query().Get("device_uuid"),
		"wallet_uuid": r.URL.Query().Get("wallet_uuid"),
		"chain_id":    r.URL.Query().Get("chain_id"),
	}

	list, total, err := rs.svc.WalletService.ListWallets(
		r.Context(),
		page,
		pageSize,
		filters,
	)
	if err != nil {
		log.Error("list wallets failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
