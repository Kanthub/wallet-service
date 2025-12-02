package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type Role struct {
	Guid       string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	RoleName   string    `gorm:"column:role_name;type:varchar(100);default:''" json:"role_name"`
	Detail     string    `gorm:"column:detail;type:varchar(255);default:''" json:"detail"`
	Status     int64     `gorm:"column:status;type:int;default:1" json:"status"`
	CreateID   int64     `gorm:"column:create_id;type:int;default:0" json:"create_id"`
	UpdateID   int64     `gorm:"column:update_id;type:int;default:0" json:"update_id"`
	CreateTime time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (Role) TableName() string {
	return "role"
}

type RoleView interface {
	GetByGuid(guid string) (*Role, error)
	GetRoleList(page, pageSize int, filters map[string]interface{}) ([]*Role, int64, error)
}

type RoleDB interface {
	RoleView

	StoreRole(role *Role) error
	StoreRoles(list []*Role) error
	UpdateRole(guid string, updates map[string]interface{}) error
}

type roleDB struct {
	gorm *gorm.DB
}

func NewRoleDB(db *gorm.DB) RoleDB {
	return &roleDB{gorm: db}
}

func (db *roleDB) StoreRole(r *Role) error {
	if err := db.gorm.Create(r).Error; err != nil {
		log.Error("StoreRole error", "err", err)
		return err
	}
	return nil
}

func (db *roleDB) StoreRoles(list []*Role) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreRoles error", "err", err)
		return err
	}
	return nil
}

func (db *roleDB) GetByGuid(guid string) (*Role, error) {
	var r Role
	if err := db.gorm.Where("guid = ?", guid).First(&r).Error; err != nil {
		log.Error("GetByGuid role error", "err", err)
		return nil, err
	}
	return &r, nil
}

func (db *roleDB) GetRoleList(page, pageSize int, filters map[string]interface{}) ([]*Role, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var list []*Role
	query := db.gorm.Model(&Role{})

	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}
		switch key {
		case "role_name", "detail":
			query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error("GetRoleList count error", "err", err)
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
		log.Error("GetRoleList list error", "err", err)
		return nil, 0, err
	}

	return list, total, nil
}

func (db *roleDB) UpdateRole(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}
	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&Role{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateRole error", "err", err)
		return err
	}
	return nil
}
