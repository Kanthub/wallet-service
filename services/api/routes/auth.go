package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"
	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) AuthManageApi() {
	r := rs.router
	r.Route("/api/v1/auth", func(r chi.Router) {

		r.Post("/create", rs.createAuth)
		r.Post("/update", rs.updateAuth)
		r.Get("/info", rs.getAuth)
		r.Get("/list", rs.listAuths)
	})
}

func (rs *Routes) createAuth(w http.ResponseWriter, r *http.Request) {
	var req service.CreateAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	auth, err := rs.svc.AuthService.CreateAuth(r.Context(), req)
	if err != nil {
		log.Error("create auth failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auth)
}

func (rs *Routes) updateAuth(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.AuthService.UpdateAuth(r.Context(), req); err != nil {
		log.Error("update auth failed", "err", err, "guid", req.Guid)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "update success",
	})
}

func (rs *Routes) getAuth(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid is required", http.StatusBadRequest)
		return
	}

	auth, err := rs.svc.AuthService.GetAuth(r.Context(), guid)
	if err != nil {
		log.Error("get auth failed", "err", err, "guid", guid)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(auth)
}

func (rs *Routes) listAuths(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	filters := map[string]interface{}{
		"auth_name": r.URL.Query().Get("auth_name"),
		"auth_url":  r.URL.Query().Get("auth_url"),
		"user_id":   r.URL.Query().Get("user_id"),
		"pid":       r.URL.Query().Get("pid"),
		"is_show":   r.URL.Query().Get("is_show"),
		"status":    r.URL.Query().Get("status"),
	}

	list, total, err := rs.svc.AuthService.ListAuths(
		r.Context(),
		page,
		pageSize,
		filters,
	)
	if err != nil {
		log.Error("list auth failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
