package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

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

	response, err := rs.svc.AdminUserLogin(loginReq)
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

	response, _ := r.svc.AdminLogout(logoutReq)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
