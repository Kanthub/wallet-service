package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/pkg/errors"

	"github.com/roothash-pay/wallet-services/cache"
	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
	model "github.com/roothash-pay/wallet-services/services/api/models/backend"
	"github.com/roothash-pay/wallet-services/services/common"
)

type AdminService interface {
	// auth
	AdminUserLogin(req model.AdminLoginRequest) (*model.AdminLoginResponse, error)
	AdminLogout(req model.AdminLogoutRequest) (*model.AdminLogoutResponse, error)

	// manage
	CreateAdmin(ctx context.Context, req CreateAdminRequest) (*backend.Admin, error)
	UpdateAdmin(ctx context.Context, req UpdateAdminRequest) error
	GetAdmin(ctx context.Context, guid string) (*backend.Admin, error)
	ListAdmins(
		ctx context.Context,
		page, pageSize int,
		filters map[string]interface{},
	) ([]*backend.Admin, int64, error)
}

type CreateAdminRequest struct {
	LoginName string `json:"login_name"`
	RealName  string `json:"real_name"`
	Password  string `json:"password"`
	RoleIDs   string `json:"role_ids"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Status    int    `json:"status"`
	CreateID  int    `json:"create_id"`
}

type UpdateAdminRequest struct {
	Guid     string                 `json:"guid"`
	Updates  map[string]interface{} `json:"updates"`
	UpdateID int                    `json:"update_id"`
}

type adminService struct {
	db           *database.DB
	siweVerifier *common.SIWEVerifier
}

func NewAdminService(db *database.DB, siweVerifier *common.SIWEVerifier) AdminService {
	return &adminService{db: db, siweVerifier: siweVerifier}
}

func (s *adminService) AdminUserLogin(
	req model.AdminLoginRequest,
) (*model.AdminLoginResponse, error) {

	if req.Username == "" || req.Password == "" {
		return &model.AdminLoginResponse{
			Success: false,
			Message: "user and password is empty",
		}, nil
	}

	user, err := s.db.BackendAdmin.GetByLoginName(req.Username)
	if err != nil {
		log.Error("query admin fail", "username", req.Username, "err", err)
		return &model.AdminLoginResponse{
			Success: false,
			Message: "user name or password is error",
		}, nil
	}

	if user.Password != req.Password {
		log.Warn("admin user name and password error", "username", req.Username)
		return &model.AdminLoginResponse{
			Success: false,
			Message: "user and password error",
		}, nil
	}

	token, err := s.siweVerifier.GenerateJWT(user.Guid, 24)
	if err != nil {
		log.Error("generate jwt fail", "err", err)
		return &model.AdminLoginResponse{
			Success: false,
			Message: "generate token failed",
		}, nil
	}

	if err := cache.Set(context.Background(), req.Username, token, 24*time.Hour); err != nil {
		log.Error("cache token fail", "err", err)
		return &model.AdminLoginResponse{
			Success: false,
			Message: "cache token failed",
		}, nil
	}

	return &model.AdminLoginResponse{
		Success: true,
		Message: "login success",
		Token:   token,
		AdminInfo: &model.AdminInfo{
			Guid:     user.Guid,
			Username: user.RealName,
		},
	}, nil
}

func (s *adminService) AdminLogout(req model.AdminLogoutRequest) (*model.AdminLogoutResponse, error) {
	if req.UserId == "" {
		log.Error("user logout failed: user_id is empty")
		return &model.AdminLogoutResponse{
			Success: false,
			Message: "user id is empty",
		}, errors.New("user id is empty")
	}

	cache.Delete(context.Background(), req.UserId)
	log.Info("user logout success, token deleted from cache", "user_id", req.UserId)

	return &model.AdminLogoutResponse{
		Success: true,
		Message: "logout success",
	}, nil
}

func (s *adminService) CreateAdmin(
	ctx context.Context,
	req CreateAdminRequest,
) (*backend.Admin, error) {

	if req.LoginName == "" || req.Password == "" {
		return nil, fmt.Errorf("login_name and password required")
	}

	admin := &backend.Admin{
		LoginName:  req.LoginName,
		RealName:   req.RealName,
		Password:   req.Password,
		RoleIDs:    req.RoleIDs,
		Phone:      req.Phone,
		Email:      req.Email,
		Status:     req.Status,
		CreateID:   req.CreateID,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := s.db.BackendAdmin.StoreAdmin(admin); err != nil {
		return nil, err
	}
	return admin, nil
}

func (s *adminService) UpdateAdmin(
	ctx context.Context,
	req UpdateAdminRequest,
) error {

	if req.Guid == "" {
		return fmt.Errorf("guid required")
	}
	if req.Updates == nil || len(req.Updates) == 0 {
		return fmt.Errorf("updates empty")
	}

	req.Updates["update_id"] = req.UpdateID
	return s.db.BackendAdmin.UpdateAdmin(req.Guid, req.Updates)
}

func (s *adminService) GetAdmin(
	ctx context.Context,
	guid string,
) (*backend.Admin, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid required")
	}
	return s.db.BackendAdmin.GetByGuid(guid)
}

func (s *adminService) ListAdmins(
	ctx context.Context,
	page, pageSize int,
	filters map[string]interface{},
) ([]*backend.Admin, int64, error) {

	return s.db.BackendAdmin.GetAdminList(page, pageSize, filters)
}
