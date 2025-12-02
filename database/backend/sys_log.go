package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type SysLog struct {
	Guid       string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	Action     string    `gorm:"column:action;type:varchar(100);default:''" json:"action"`
	Remark     string    `gorm:"column:remark;type:varchar(100);default:''" json:"remark"`
	Admin      string    `gorm:"column:admin;type:varchar(30);default:''" json:"admin"`
	IP         string    `gorm:"column:ip;type:varchar(30);default:''" json:"ip"`
	Cate       int64     `gorm:"column:cate;type:smallint;default:0" json:"cate"`
	Status     int64     `gorm:"column:status;type:smallint;default:-1" json:"status"`
	Asset      string    `gorm:"column:asset;type:varchar(255);default:''" json:"asset"`
	Before     string    `gorm:"column:before;type:varchar(255);default:''" json:"before"`
	After      string    `gorm:"column:after;type:varchar(255);default:''" json:"after"`
	UserID     int64     `gorm:"column:user_id;type:bigint;default:0" json:"user_id"`
	OrderNo    string    `gorm:"column:order_number;type:varchar(64);default:''" json:"order_number"`
	Op         int64     `gorm:"column:op;type:smallint;default:-1" json:"op"`
	CreateTime time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (SysLog) TableName() string {
	return "sys_log"
}

type SysLogView interface {
	GetByGuid(guid string) (*SysLog, error)
	GetSysLogList(page, pageSize int, filters map[string]interface{}) ([]*SysLog, int64, error)
}

type SysLogDB interface {
	SysLogView

	StoreSysLog(logItem *SysLog) error
	StoreSysLogs(list []*SysLog) error
	UpdateSysLog(guid string, updates map[string]interface{}) error
}

type sysLogDB struct {
	gorm *gorm.DB
}

func NewSysLogDB(db *gorm.DB) SysLogDB {
	return &sysLogDB{gorm: db}
}

func (db *sysLogDB) StoreSysLog(item *SysLog) error {
	if err := db.gorm.Create(item).Error; err != nil {
		log.Error("StoreSysLog error", "err", err)
		return err
	}
	return nil
}

func (db *sysLogDB) StoreSysLogs(list []*SysLog) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreSysLogs error", "err", err)
		return err
	}
	return nil
}

func (db *sysLogDB) GetByGuid(guid string) (*SysLog, error) {
	var item SysLog
	if err := db.gorm.Where("guid = ?", guid).First(&item).Error; err != nil {
		log.Error("GetByGuid sys_log error", "err", err)
		return nil, err
	}
	return &item, nil
}

func (db *sysLogDB) GetSysLogList(page, pageSize int, filters map[string]interface{}) ([]*SysLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var list []*SysLog
	query := db.gorm.Model(&SysLog{})

	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}
		switch key {
		case "action", "remark", "asset", "order_number":
			query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error("GetSysLogList count error", "err", err)
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
		log.Error("GetSysLogList list error", "err", err)
		return nil, 0, err
	}

	return list, total, nil
}

func (db *sysLogDB) UpdateSysLog(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}
	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&SysLog{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateSysLog error", "err", err)
		return err
	}
	return nil
}
