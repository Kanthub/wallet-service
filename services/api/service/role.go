package service

import (
	"context"
	"fmt"
	"time"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type CreateRoleRequest struct {
	RoleName string `json:"role_name"`
	Detail   string `json:"detail"`
	Status   int64  `json:"status"`
	CreateID int64  `json:"create_id"`
}

type UpdateRoleRequest struct {
	Guid     string                 `json:"guid"`
	Updates  map[string]interface{} `json:"updates"`
	UpdateID int64                  `json:"update_id"`
}

type RoleService interface {
	CreateRole(ctx context.Context, req CreateRoleRequest) (*backend.Role, error)
	UpdateRole(ctx context.Context, req UpdateRoleRequest) error
	GetRole(ctx context.Context, guid string) (*backend.Role, error)
	ListRoles(
		ctx context.Context,
		page, pageSize int,
		filters map[string]interface{},
	) ([]*backend.Role, int64, error)
}

type roleService struct {
	db *database.DB
}

func NewRoleService(db *database.DB) RoleService {
	return &roleService{db: db}
}

func (s *roleService) CreateRole(
	ctx context.Context,
	req CreateRoleRequest,
) (*backend.Role, error) {

	if req.RoleName == "" {
		return nil, fmt.Errorf("role_name is required")
	}

	role := &backend.Role{
		RoleName:   req.RoleName,
		Detail:     req.Detail,
		Status:     req.Status,
		CreateID:   req.CreateID,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	if err := s.db.BackendRole.StoreRole(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *roleService) UpdateRole(
	ctx context.Context,
	req UpdateRoleRequest,
) error {

	if req.Guid == "" {
		return fmt.Errorf("guid is required")
	}
	if len(req.Updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	req.Updates["update_id"] = req.UpdateID
	return s.db.BackendRole.UpdateRole(req.Guid, req.Updates)
}

func (s *roleService) GetRole(
	ctx context.Context,
	guid string,
) (*backend.Role, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid is required")
	}
	return s.db.BackendRole.GetByGuid(guid)
}

func (s *roleService) ListRoles(
	ctx context.Context,
	page, pageSize int,
	filters map[string]interface{},
) ([]*backend.Role, int64, error) {

	return s.db.BackendRole.GetRoleList(page, pageSize, filters)
}
