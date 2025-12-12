package service

import (
	"context"
	"fmt"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type SysLogService interface {
	// 写日志（给 middleware / handler 用）
	WriteLog(ctx context.Context, item *backend.SysLog) error

	// 查日志（给后台页面用）
	ListLogs(
		ctx context.Context,
		page, pageSize int,
		filters map[string]interface{},
	) ([]*backend.SysLog, int64, error)

	// 查单条（可选）
	GetLog(ctx context.Context, guid string) (*backend.SysLog, error)
}

type sysLogService struct {
	db *database.DB
}

func NewSysLogService(db *database.DB) SysLogService {
	return &sysLogService{db: db}
}

func (s *sysLogService) WriteLog(
	ctx context.Context,
	item *backend.SysLog,
) error {

	if item == nil {
		return fmt.Errorf("syslog item is nil")
	}
	if item.Action == "" {
		return fmt.Errorf("action is required")
	}

	return s.db.BackendSysLog.StoreSysLog(item)
}

func (s *sysLogService) ListLogs(
	ctx context.Context,
	page, pageSize int,
	filters map[string]interface{},
) ([]*backend.SysLog, int64, error) {

	return s.db.BackendSysLog.GetSysLogList(page, pageSize, filters)
}

func (s *sysLogService) GetLog(
	ctx context.Context,
	guid string,
) (*backend.SysLog, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid is required")
	}
	return s.db.BackendSysLog.GetByGuid(guid)
}
