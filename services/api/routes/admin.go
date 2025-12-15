package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"
	"github.com/roothash-pay/wallet-services/services/api/models/backend"
	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) AdminManageApi() {
	r := rs.router
	r.Route("/api/v1/admin", func(r chi.Router) {

		r.Post("/create", rs.createAdmin)
		r.Post("/update", rs.updateAdmin)
		r.Get("/info", rs.getAdmin)
		r.Get("/list", rs.listAdmins)

		r.Post("/login", rs.AdminLoginHandler)
		r.Post("/logout", rs.AdminLogoutHandler)
	})
}

func (rs *Routes) createAdmin(w http.ResponseWriter, r *http.Request) {
	var req service.CreateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	admin, err := rs.svc.AdminService.CreateAdmin(r.Context(), req)
	if err != nil {
		log.Error("create admin failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(admin)
}

func (rs *Routes) updateAdmin(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.AdminService.UpdateAdmin(r.Context(), req); err != nil {
		log.Error("update admin failed", "err", err, "guid", req.Guid)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "update success",
	})
}

func (rs *Routes) getAdmin(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid is required", http.StatusBadRequest)
		return
	}

	admin, err := rs.svc.AdminService.GetAdmin(r.Context(), guid)
	if err != nil {
		log.Error("get admin failed", "err", err, "guid", guid)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(admin)
}

func (rs *Routes) listAdmins(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	filters := map[string]interface{}{
		"login_name": r.URL.Query().Get("login_name"),
		"real_name":  r.URL.Query().Get("real_name"),
		"email":      r.URL.Query().Get("email"),
		"phone":      r.URL.Query().Get("phone"),
		"status":     r.URL.Query().Get("status"),
	}

	list, total, err := rs.svc.AdminService.ListAdmins(
		r.Context(),
		page,
		pageSize,
		filters,
	)
	if err != nil {
		log.Error("list admins failed", "err", err)
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

func (rs *Routes) AdminLoginHandler(w http.ResponseWriter, req *http.Request) {
	var loginReq backend.AdminLoginRequest
	if err := json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(backend.AdminLoginResponse{
			Success: false,
			Message: "invalid params",
		})
		log.Error("parse params fail", "err", err)
		return
	}

	response, err := rs.svc.AdminService.AdminUserLogin(loginReq)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(backend.AdminLoginResponse{
			Success: false,
			Message: "login fail, try again later",
		})
		log.Error("admin user login fail", "err", err)
		return
	}

	statusCode := http.StatusOK
	if !response.Success {
		statusCode = http.StatusBadRequest
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (r *Routes) AdminLogoutHandler(w http.ResponseWriter, req *http.Request) {
	var logoutReq backend.AdminLogoutRequest
	if err := json.NewDecoder(req.Body).Decode(&logoutReq); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(backend.AdminLoginResponse{
			Success: false,
			Message: "invalid request body",
		})
		return
	}

	response, _ := r.svc.AdminService.AdminLogout(logoutReq)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
