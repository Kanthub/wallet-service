package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) WalletAddressNoteApi() {
	r := rs.router
	r.Route("/api/v1/wallet-address-note", func(r chi.Router) {

		r.Post("/create", rs.createWalletAddressNote)
		r.Post("/update", rs.updateWalletAddressNote)
		r.Get("/info", rs.getWalletAddressNote)
		r.Get("/list", rs.getWalletAddressNoteList)
	})
}

func (rs *Routes) createWalletAddressNote(w http.ResponseWriter, r *http.Request) {
	var req service.CreateWalletAddressNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.WalletAddressNoteService.CreateWalletAddressNote(
		r.Context(),
		req,
	)
	if err != nil {
		log.Error("create wallet_address_note failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) updateWalletAddressNote(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateWalletAddressNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.WalletAddressNoteService.UpdateWalletAddressNote(
		r.Context(),
		req,
	); err != nil {
		log.Error("update wallet_address_note failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "update success",
	})
}

func (rs *Routes) getWalletAddressNote(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.WalletAddressNoteService.GetWalletAddressNote(
		r.Context(),
		guid,
	)
	if err != nil {
		log.Error("get wallet_address_note failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) getWalletAddressNoteList(w http.ResponseWriter, r *http.Request) {
	deviceUUID := r.URL.Query().Get("device_uuid")
	if deviceUUID == "" {
		http.Error(w, "device_uuid required", http.StatusBadRequest)
		return
	}

	list, err := rs.svc.WalletAddressNoteService.GetByDeviceUUID(
		r.Context(),
		deviceUUID,
	)
	if err != nil {
		log.Error("get wallet_address_note list failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(list)
}
