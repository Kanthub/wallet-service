package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/common/httputil"
	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) AddressAssetApi() {
	r := rs.router
	r.Route("/api/v1/addressAsset", func(r chi.Router) {
		r.Post("/upsert", rs.upsertAsset)
		r.Get("/list", rs.listAsset)
		r.Get("/info", rs.getAsset)
	})
}

func (rs *Routes) upsertAsset(w http.ResponseWriter, r *http.Request) {
	var req service.UpsertAddressAssetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, "invalid request")
		return
	}

	asset, err := rs.svc.AddressAssetService.UpsertAddressAsset(r.Context(), req)
	if err != nil {
		httputil.WriteError(w, err.Error())
		return
	}

	httputil.WriteSuccess(w, asset)
}

func (rs *Routes) listAsset(w http.ResponseWriter, r *http.Request) {
	addressUUID := r.URL.Query().Get("address_uuid")
	if addressUUID == "" {
		httputil.WriteError(w, "address_uuid required")
		return
	}

	list, err := rs.svc.AddressAssetService.ListAddressAssetsByAddress(r.Context(), addressUUID)
	if err != nil {
		httputil.WriteError(w, err.Error())
		return
	}

	httputil.WriteSuccess(w, list)
}

func (rs *Routes) getAsset(w http.ResponseWriter, r *http.Request) {
	addressUUID := r.URL.Query().Get("address_uuid")
	tokenID := r.URL.Query().Get("token_id")

	if addressUUID == "" || tokenID == "" {
		httputil.WriteError(w, "missing parameters")
		return
	}

	asset, err := rs.svc.AddressAssetService.GetAddressAsset(r.Context(), addressUUID, tokenID)
	if err != nil {
		httputil.WriteError(w, err.Error())
		return
	}

	httputil.WriteSuccess(w, asset)
}
