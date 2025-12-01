package backend

import (
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/log"
)

type Role struct {
	ID         int64  `gorm:"primaryKey;column:id" json:"id"`
	RoleName   string `gorm:"type:varchar(100);default:''" json:"role_name"` // 角色名称
	Detail     string `gorm:"type:varchar(255);default:''" json:"detail"`    // 角色描述/说明
	Status     int    `gorm:"type:int;default:1" json:"status"`              // 状态(1启用;0禁用)
	CreateID   int    `gorm:"type:int;default:0" json:"create_id"`           // 创建人ID
	UpdateID   int    `gorm:"type:int;default:0" json:"update_id"`           // 修改人ID
	CreateTime int64  `gorm:"type:bigint;default:0" json:"create_time"`      // 创建时间(Unix时间戳)
	UpdateTime int64  `gorm:"type:bigint;default:0" json:"update_time"`      // 更新时间(Unix时间戳)
}

func (Role) TableName() string {
	return "role"
}

type RoleView interface {
	GetByID(id int64) (*Role, error)
	GetByStatus(status int) ([]*Role, error)
}

type RoleDB interface {
	RoleView

	StoreRole(role *Role) error
	StoreRoles(roles []*Role) error
}

type roleDB struct {
	gorm *gorm.DB
}

func NewRoleDB(db *gorm.DB) RoleDB {
	return &roleDB{gorm: db}
}

func (db *roleDB) StoreRole(role *Role) error {
	if err := db.gorm.Create(role).Error; err != nil {
		log.Error("StoreRole error:", err)
		return err
	}
	return nil
}

func (db *roleDB) StoreRoles(roles []*Role) error {
	if err := db.gorm.CreateInBatches(roles, len(roles)).Error; err != nil {
		log.Error("StoreRoles error:", err)
		return err
	}
	return nil
}

func (db *roleDB) GetByID(id int64) (*Role, error) {
	var role Role
	if err := db.gorm.First(&role, id).Error; err != nil {
		log.Error("GetByID error:", err)
		return nil, err
	}
	return &role, nil
}

func (db *roleDB) GetByStatus(status int) ([]*Role, error) {
	var list []*Role
	if err := db.gorm.Where("status = ?", status).Find(&list).Error; err != nil {
		log.Error("GetByStatus error:", err)
		return nil, err
	}
	return list, nil
}
