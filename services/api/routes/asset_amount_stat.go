package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) AssetAmountStatApi() {
	r := rs.router
	r.Route("/api/v1/asset-amount-stat", func(r chi.Router) {

		r.Post("/create", rs.createAssetAmountStat)
		r.Post("/update", rs.updateAssetAmountStat)
		r.Get("/by-asset-date", rs.getAssetAmountByDate)
	})
}

func (rs *Routes) createAssetAmountStat(w http.ResponseWriter, r *http.Request) {
	var req service.CreateAssetAmountStatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.AssetAmountStatService.CreateStat(r.Context(), req)
	if err != nil {
		log.Error("create asset_amount_stat failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) updateAssetAmountStat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Guid    string                 `json:"guid"`
		Updates map[string]interface{} `json:"updates"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.AssetAmountStatService.UpdateStat(r.Context(), req.Guid, req.Updates); err != nil {
		log.Error("update asset_amount_stat failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func (rs *Routes) getAssetAmountByDate(w http.ResponseWriter, r *http.Request) {
	assetUUID := r.URL.Query().Get("asset_uuid")
	date := r.URL.Query().Get("date")

	if assetUUID == "" || date == "" {
		http.Error(w, "asset_uuid and date required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.AssetAmountStatService.GetByAssetAndDate(r.Context(), assetUUID, date)
	if err != nil {
		log.Error("get asset_amount_stat failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}
