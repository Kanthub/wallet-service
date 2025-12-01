package service

import (
	"context"
	"github.com/pkg/errors"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/cache"
	"github.com/roothash-pay/wallet-services/services/api/models/backend"
)

func (h *HandlerSvc) AdminUserLogin(req backend.AdminLoginRequest) (*backend.AdminLoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return &backend.AdminLoginResponse{
			Success: false,
			Message: "user and password is empty",
		}, nil
	}

	user, err := h.backendAdminDB.GetByLoginName(req.Username)
	if err != nil {
		log.Error("query admin fail", "username", req.Username, "err", err)
		return &backend.AdminLoginResponse{
			Success: false,
			Message: "user name or password is error",
		}, nil
	}

	if user.Password != req.Password {
		log.Warn("admin user name and password error", "username", req.Username)
		return &backend.AdminLoginResponse{
			Success: false,
			Message: "user and password error",
		}, nil
	}

	token, err := h.siweVerifier.GenerateJWT(strconv.FormatInt(user.ID, 10), 24) // 24小时有效期
	if err != nil {
		log.Error("Failed to generate JWT", "err", err)
		return &backend.AdminLoginResponse{
			Success: false,
			Message: "failed to generate authentication token",
		}, nil
	}

	if err := cache.Set(context.Background(), req.Username, token, 24*time.Hour); err != nil {
		log.Error("Failed to save token to cache", "err", err)
		return &backend.AdminLoginResponse{
			Success: false,
			Message: "failed to save authentication token",
		}, nil
	}

	return &backend.AdminLoginResponse{
		Success: true,
		Message: "login success",
		Token:   token,
		AdminInfo: &backend.AdminInfo{
			ID:       uint64(user.ID),
			Username: user.RealName,
		},
	}, nil
}

func (h *HandlerSvc) AdminLogout(req backend.AdminLogoutRequest) (*backend.AdminLogoutResponse, error) {
	if req.UserId == "" {
		log.Error("user logout failed: user_id is empty")
		return &backend.AdminLogoutResponse{
			Success: false,
			Message: "user id is empty",
		}, errors.New("user id is empty")
	}

	cache.Delete(context.Background(), req.UserId)
	log.Info("user logout success, token deleted from cache", "user_id", req.UserId)

	return &backend.AdminLogoutResponse{
		Success: true,
		Message: "logout success",
	}, nil
}
