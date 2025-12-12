package service

import (
	"context"
	"fmt"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type BindRoleAuthRequest struct {
	RoleID  int64   `json:"role_id"`
	AuthIDs []int64 `json:"auth_ids"`
}

type BindAuthRoleRequest struct {
	AuthID  int64   `json:"auth_id"`
	RoleIDs []int64 `json:"role_ids"`
}

type RoleAuthService interface {
	BindAuthsToRole(ctx context.Context, req BindRoleAuthRequest) error
	BindRolesToAuth(ctx context.Context, req BindAuthRoleRequest) error

	GetAuthsByRole(ctx context.Context, roleID int64) ([]*backend.RoleAuth, error)
	GetRolesByAuth(ctx context.Context, authID int64) ([]*backend.RoleAuth, error)
}

type roleAuthService struct {
	db *database.DB
}

func NewRoleAuthService(db *database.DB) RoleAuthService {
	return &roleAuthService{db: db}
}

func (s *roleAuthService) BindAuthsToRole(
	ctx context.Context,
	req BindRoleAuthRequest,
) error {

	if req.RoleID <= 0 {
		return fmt.Errorf("invalid role_id")
	}

	// 先清空该角色下的所有权限
	if err := s.db.BackendRoleAuth.DeleteByRoleID(req.RoleID); err != nil {
		return err
	}

	if len(req.AuthIDs) == 0 {
		return nil
	}

	list := make([]*backend.RoleAuth, 0, len(req.AuthIDs))
	for _, authID := range req.AuthIDs {
		if authID <= 0 {
			continue
		}
		list = append(list, &backend.RoleAuth{
			RoleID: req.RoleID,
			AuthID: authID,
		})
	}

	if len(list) == 0 {
		return nil
	}

	return s.db.BackendRoleAuth.StoreRoleAuths(list)
}

func (s *roleAuthService) BindRolesToAuth(
	ctx context.Context,
	req BindAuthRoleRequest,
) error {

	if req.AuthID <= 0 {
		return fmt.Errorf("invalid auth_id")
	}

	if err := s.db.BackendRoleAuth.DeleteByAuthID(req.AuthID); err != nil {
		return err
	}

	if len(req.RoleIDs) == 0 {
		return nil
	}

	list := make([]*backend.RoleAuth, 0, len(req.RoleIDs))
	for _, roleID := range req.RoleIDs {
		if roleID <= 0 {
			continue
		}
		list = append(list, &backend.RoleAuth{
			AuthID: req.AuthID,
			RoleID: roleID,
		})
	}

	if len(list) == 0 {
		return nil
	}

	return s.db.BackendRoleAuth.StoreRoleAuths(list)
}

func (s *roleAuthService) GetAuthsByRole(
	ctx context.Context,
	roleID int64,
) ([]*backend.RoleAuth, error) {

	if roleID <= 0 {
		return nil, fmt.Errorf("invalid role_id")
	}
	return s.db.BackendRoleAuth.GetByRoleID(roleID)
}

func (s *roleAuthService) GetRolesByAuth(
	ctx context.Context,
	authID int64,
) ([]*backend.RoleAuth, error) {

	if authID <= 0 {
		return nil, fmt.Errorf("invalid auth_id")
	}
	return s.db.BackendRoleAuth.GetByAuthID(authID)
}
