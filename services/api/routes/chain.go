package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"
	"github.com/roothash-pay/wallet-services/database/backend"
)

func (rs *Routes) ChainApi() {
	r := rs.router
	r.Route("/api/v1/chain", func(r chi.Router) {

		// 管理后台
		r.Post("/create", rs.createChain)
		r.Post("/update", rs.updateChain)

		// 查询
		r.Get("/info", rs.getChain)
		r.Get("/by_chain_id", rs.getChainByChainID)
		r.Get("/by_name", rs.getChainByName)
		r.Get("/list", rs.listChains)

		// 前端 / 公共
		r.Get("/all", rs.listAllChains)
	})
}

func (rs *Routes) createChain(w http.ResponseWriter, r *http.Request) {
	var c backend.Chain
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	chain, err := rs.svc.ChainService.CreateChain(r.Context(), &c)
	if err != nil {
		log.Error("create chain failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(chain)
}

func (rs *Routes) updateChain(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Guid    string                 `json:"guid"`
		Updates map[string]interface{} `json:"updates"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.ChainService.UpdateChain(r.Context(), req.Guid, req.Updates); err != nil {
		log.Error("update chain failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
	})
}

func (rs *Routes) getChain(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid required", http.StatusBadRequest)
		return
	}

	chain, err := rs.svc.ChainService.GetChain(r.Context(), guid)
	if err != nil {
		log.Error("get chain failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(chain)
}

func (rs *Routes) getChainByChainID(w http.ResponseWriter, r *http.Request) {
	chainID := r.URL.Query().Get("chain_id")
	if chainID == "" {
		http.Error(w, "chain_id required", http.StatusBadRequest)
		return
	}

	chain, err := rs.svc.ChainService.GetChainByChainID(r.Context(), chainID)
	if err != nil {
		log.Error("get chain by chain_id failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(chain)
}

func (rs *Routes) getChainByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("chain_name")
	if name == "" {
		http.Error(w, "chain_name required", http.StatusBadRequest)
		return
	}

	chain, err := rs.svc.ChainService.GetChainByName(r.Context(), name)
	if err != nil {
		log.Error("get chain by name failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(chain)
}

func (rs *Routes) listChains(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	filters := map[string]interface{}{
		"chain_name": r.URL.Query().Get("chain_name"),
		"chain_mark": r.URL.Query().Get("chain_mark"),
		"network":    r.URL.Query().Get("network"),
		"is_enabled": r.URL.Query().Get("is_enabled"),
	}

	list, total, err := rs.svc.ChainService.ListChains(
		r.Context(),
		page,
		pageSize,
		filters,
	)
	if err != nil {
		log.Error("list chains failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (rs *Routes) listAllChains(w http.ResponseWriter, r *http.Request) {
	list, err := rs.svc.ChainService.ListAllChains(r.Context())
	if err != nil {
		log.Error("list all chains failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(list)
}
