package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) WalletTxRecordApi() {
	r := rs.router
	r.Route("/api/v1/wallet-tx", func(r chi.Router) {

		r.Post("/create", rs.createWalletTx)
		r.Post("/update", rs.updateWalletTx)
		r.Get("/info", rs.getWalletTx)
		r.Get("/by-operation", rs.getWalletTxByOperation)
	})
}

func (rs *Routes) createWalletTx(w http.ResponseWriter, r *http.Request) {
	var req service.CreateWalletTxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.WalletTxRecordService.CreateWalletTx(r.Context(), req)
	if err != nil {
		log.Error("create wallet_tx_record failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) updateWalletTx(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateWalletTxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.WalletTxRecordService.UpdateWalletTx(r.Context(), req); err != nil {
		log.Error("update wallet_tx_record failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "update success",
	})
}

func (rs *Routes) getWalletTx(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.WalletTxRecordService.GetWalletTx(r.Context(), guid)
	if err != nil {
		log.Error("get wallet_tx_record failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) getWalletTxByOperation(w http.ResponseWriter, r *http.Request) {
	operationID := r.URL.Query().Get("operation_id")
	if operationID == "" {
		http.Error(w, "operation_id required", http.StatusBadRequest)
		return
	}

	list, err := rs.svc.WalletTxRecordService.GetByOperationID(
		r.Context(),
		operationID,
	)
	if err != nil {
		log.Error("get wallet_tx_record by operation failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(list)
}
