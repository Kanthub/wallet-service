package service

import (
	"context"
	"fmt"
	"time"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type CreateAuthRequest struct {
	AuthName string `json:"auth_name"`
	AuthURL  string `json:"auth_url"`
	UserID   int64  `json:"user_id"`
	PID      int64  `json:"pid"`
	Sort     int64  `json:"sort"`
	Icon     string `json:"icon"`
	IsShow   int64  `json:"is_show"`
	Status   int64  `json:"status"`
	CreateID int64  `json:"create_id"`
}

type UpdateAuthRequest struct {
	Guid     string                 `json:"guid"`
	Updates  map[string]interface{} `json:"updates"`
	UpdateID int64                  `json:"update_id"`
}

type AuthService interface {
	CreateAuth(ctx context.Context, req CreateAuthRequest) (*backend.Auth, error)
	UpdateAuth(ctx context.Context, req UpdateAuthRequest) error
	GetAuth(ctx context.Context, guid string) (*backend.Auth, error)
	ListAuths(
		ctx context.Context,
		page, pageSize int,
		filters map[string]interface{},
	) ([]*backend.Auth, int64, error)
}

type authService struct {
	db *database.DB
}

func NewAuthService(db *database.DB) AuthService {
	return &authService{db: db}
}

func (s *authService) CreateAuth(
	ctx context.Context,
	req CreateAuthRequest,
) (*backend.Auth, error) {

	if req.AuthName == "" {
		return nil, fmt.Errorf("auth_name is required")
	}

	auth := &backend.Auth{
		AuthName:   req.AuthName,
		AuthURL:    req.AuthURL,
		UserID:     req.UserID,
		PID:        req.PID,
		Sort:       req.Sort,
		Icon:       req.Icon,
		IsShow:     req.IsShow,
		Status:     req.Status,
		CreateID:   req.CreateID,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := s.db.BackendAuth.StoreAuth(auth); err != nil {
		return nil, err
	}
	return auth, nil
}

func (s *authService) UpdateAuth(
	ctx context.Context,
	req UpdateAuthRequest,
) error {

	if req.Guid == "" {
		return fmt.Errorf("guid is required")
	}
	if len(req.Updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	req.Updates["update_id"] = req.UpdateID
	return s.db.BackendAuth.UpdateAuth(req.Guid, req.Updates)
}

func (s *authService) GetAuth(
	ctx context.Context,
	guid string,
) (*backend.Auth, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid is required")
	}
	return s.db.BackendAuth.GetByGuid(guid)
}

func (s *authService) ListAuths(
	ctx context.Context,
	page, pageSize int,
	filters map[string]interface{},
) ([]*backend.Auth, int64, error) {

	return s.db.BackendAuth.GetAuthList(page, pageSize, filters)
}
