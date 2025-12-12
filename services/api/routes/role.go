package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"
	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) RoleManageApi() {
	r := rs.router
	r.Route("/api/v1/role", func(r chi.Router) {

		r.Post("/create", rs.createRole)
		r.Post("/update", rs.updateRole)
		r.Get("/info", rs.getRole)
		r.Get("/list", rs.listRoles)
	})
}

func (rs *Routes) createRole(w http.ResponseWriter, r *http.Request) {
	var req service.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	role, err := rs.svc.RoleService.CreateRole(r.Context(), req)
	if err != nil {
		log.Error("create role failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func (rs *Routes) updateRole(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.RoleService.UpdateRole(r.Context(), req); err != nil {
		log.Error("update role failed", "err", err, "guid", req.Guid)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "update success",
	})
}

func (rs *Routes) getRole(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid is required", http.StatusBadRequest)
		return
	}

	role, err := rs.svc.RoleService.GetRole(r.Context(), guid)
	if err != nil {
		log.Error("get role failed", "err", err, "guid", guid)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func (rs *Routes) listRoles(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	filters := map[string]interface{}{
		"role_name": r.URL.Query().Get("role_name"),
		"detail":    r.URL.Query().Get("detail"),
		"status":    r.URL.Query().Get("status"),
	}

	list, total, err := rs.svc.RoleService.ListRoles(
		r.Context(),
		page,
		pageSize,
		filters,
	)
	if err != nil {
		log.Error("list roles failed", "err", err)
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
