package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) TokenApi() {
	r := rs.router
	r.Route("/api/v1/token", func(r chi.Router) {

		r.Post("/create", rs.createToken)
		r.Post("/update", rs.updateToken)
		r.Get("/info", rs.getToken)
		r.Get("/list", rs.listTokens)
	})
}

func (rs *Routes) createToken(w http.ResponseWriter, r *http.Request) {
	var req service.CreateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := rs.svc.TokenService.CreateToken(r.Context(), req)
	if err != nil {
		log.Error("create token failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(token)
}

func (rs *Routes) updateToken(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.TokenService.UpdateToken(r.Context(), req); err != nil {
		log.Error("update token failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "update success",
	})
}

func (rs *Routes) getToken(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid required", http.StatusBadRequest)
		return
	}

	token, err := rs.svc.TokenService.GetToken(r.Context(), guid)
	if err != nil {
		log.Error("get token failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(token)
}

func (rs *Routes) listTokens(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	filters := map[string]interface{}{
		"token_name":   r.URL.Query().Get("token_name"),
		"token_symbol": r.URL.Query().Get("token_symbol"),
		"chain_id":     r.URL.Query().Get("chain_id"),
		"is_hot":       r.URL.Query().Get("is_hot"),
	}

	list, total, err := rs.svc.TokenService.ListTokens(
		r.Context(),
		page,
		pageSize,
		filters,
	)
	if err != nil {
		log.Error("list tokens failed", "err", err)
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
