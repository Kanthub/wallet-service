package backend

import (
	"gorm.io/gorm"
	"time"
)

type SysLog struct {
	ID          int64     `gorm:"primaryKey;column:id" json:"id"`
	Action      string    `gorm:"type:varchar(100);default:''" json:"action"` // 路径
	Desc        string    `gorm:"type:varchar(100);default:''" json:"desc"`   // 描述
	Admin       string    `gorm:"type:varchar(30);default:''" json:"admin"`   // 管理员
	IP          string    `gorm:"type:varchar(30);default:''" json:"ip"`      // IP
	Cate        int16     `gorm:"type:smallint;default:0" json:"cate"`        // 类型
	Status      int16     `gorm:"type:smallint;default:-1" json:"status"`     // 登陆状态
	Asset       string    `gorm:"type:varchar(255);default:''" json:"asset"`  // 币种
	Before      string    `gorm:"type:varchar(255);default:''" json:"before"` // 修改前
	After       string    `gorm:"type:varchar(255);default:''" json:"after"`  // 修改后
	UserID      int64     `gorm:"type:bigint;default:0" json:"user_id"`
	OrderNumber string    `gorm:"type:varchar(64);default:''" json:"order_number"`
	Op          int16     `gorm:"type:smallint;default:-1" json:"op"` // 操作类型
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SysLog) TableName() string {
	return "sys_log"
}

type SysLogView interface {
	GetLogsByCate(cate int16) ([]*SysLog, error)
	GetLogsByStatus(status int16) ([]*SysLog, error)
	GetLogByOrderNumber(orderNumber string) (*SysLog, error)
}

type SysLogDB interface {
	SysLogView

	StoreLog(*SysLog) error
	StoreLogs([]*SysLog) error
}

type sysLogDB struct {
	gorm *gorm.DB
}

func (s *sysLogDB) GetLogsByCate(cate int16) ([]*SysLog, error) {
	panic("implement me")
}

func (s *sysLogDB) GetLogsByStatus(status int16) ([]*SysLog, error) {
	panic("implement me")
}

func (s *sysLogDB) GetLogByOrderNumber(orderNumber string) (*SysLog, error) {
	panic("implement me")
}

func (s *sysLogDB) StoreLog(log *SysLog) error {
	panic("implement me")
}

func (s *sysLogDB) StoreLogs(logs []*SysLog) error {
	panic("implement me")
}

func NewSysLogDB(db *gorm.DB) SysLogDB {
	return &sysLogDB{gorm: db}
}
