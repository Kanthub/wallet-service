package backend

import (
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/log"
)

type RoleAuth struct {
	AuthID int   `gorm:"primaryKey;column:auth_id" json:"auth_id"`
	RoleID int64 `gorm:"primaryKey;column:role_id" json:"role_id"`
}

func (RoleAuth) TableName() string {
	return "role_auth"
}

type RoleAuthView interface {
	GetByRoleID(roleID int64) ([]*RoleAuth, error)
	GetByAuthID(authID int) ([]*RoleAuth, error)
}

type RoleAuthDB interface {
	RoleAuthView

	StoreRoleAuth(roleAuth *RoleAuth) error
	StoreRoleAuths(list []*RoleAuth) error
}

type roleAuthDB struct {
	gorm *gorm.DB
}

func NewRoleAuthDB(db *gorm.DB) RoleAuthDB {
	return &roleAuthDB{gorm: db}
}

func (db *roleAuthDB) StoreRoleAuth(roleAuth *RoleAuth) error {
	if err := db.gorm.Create(roleAuth).Error; err != nil {
		log.Error("StoreRoleAuth error:", err)
		return err
	}
	return nil
}

func (db *roleAuthDB) StoreRoleAuths(list []*RoleAuth) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreRoleAuths error:", err)
		return err
	}
	return nil
}

func (db *roleAuthDB) GetByRoleID(roleID int64) ([]*RoleAuth, error) {
	var list []*RoleAuth
	if err := db.gorm.Where("role_id = ?", roleID).Find(&list).Error; err != nil {
		log.Error("GetByRoleID error:", err)
		return nil, err
	}
	return list, nil
}

func (db *roleAuthDB) GetByAuthID(authID int) ([]*RoleAuth, error) {
	var list []*RoleAuth
	if err := db.gorm.Where("auth_id = ?", authID).Find(&list).Error; err != nil {
		log.Error("GetByAuthID error:", err)
		return nil, err
	}
	return list, nil
}
