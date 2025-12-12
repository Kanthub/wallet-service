package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"
	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) RoleAuthManageApi() {
	r := rs.router
	r.Route("/api/v1/role-auth", func(r chi.Router) {

		r.Post("/bind-auths", rs.bindAuthsToRole)
		r.Post("/bind-roles", rs.bindRolesToAuth)

		r.Get("/auths", rs.getAuthsByRole)
		r.Get("/roles", rs.getRolesByAuth)
	})
}

func (rs *Routes) bindAuthsToRole(w http.ResponseWriter, r *http.Request) {
	var req service.BindRoleAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.RoleAuthService.BindAuthsToRole(r.Context(), req); err != nil {
		log.Error("bind auths to role failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "bind success",
	})
}

func (rs *Routes) bindRolesToAuth(w http.ResponseWriter, r *http.Request) {
	var req service.BindAuthRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.RoleAuthService.BindRolesToAuth(r.Context(), req); err != nil {
		log.Error("bind roles to auth failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "bind success",
	})
}

func (rs *Routes) getAuthsByRole(w http.ResponseWriter, r *http.Request) {
	roleID, _ := strconv.ParseInt(r.URL.Query().Get("role_id"), 10, 64)
	if roleID <= 0 {
		http.Error(w, "invalid role_id", http.StatusBadRequest)
		return
	}

	list, err := rs.svc.RoleAuthService.GetAuthsByRole(r.Context(), roleID)
	if err != nil {
		log.Error("get auths by role failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (rs *Routes) getRolesByAuth(w http.ResponseWriter, r *http.Request) {
	authID, _ := strconv.ParseInt(r.URL.Query().Get("auth_id"), 10, 64)
	if authID <= 0 {
		http.Error(w, "invalid auth_id", http.StatusBadRequest)
		return
	}

	list, err := rs.svc.RoleAuthService.GetRolesByAuth(r.Context(), authID)
	if err != nil {
		log.Error("get roles by auth failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
