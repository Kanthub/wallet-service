package backend

import (
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type RoleAuth struct {
	AuthID int64 `gorm:"column:auth_id;primaryKey" json:"auth_id"`
	RoleID int64 `gorm:"column:role_id;primaryKey" json:"role_id"`
}

func (RoleAuth) TableName() string {
	return "role_auth"
}

type RoleAuthView interface {
	GetByRoleID(roleID int64) ([]*RoleAuth, error)
	GetByAuthID(authID int64) ([]*RoleAuth, error)
}

type RoleAuthDB interface {
	RoleAuthView

	StoreRoleAuth(item *RoleAuth) error
	StoreRoleAuths(list []*RoleAuth) error
	DeleteByRoleID(roleID int64) error
	DeleteByAuthID(authID int64) error
}

type roleAuthDB struct {
	gorm *gorm.DB
}

func NewRoleAuthDB(db *gorm.DB) RoleAuthDB {
	return &roleAuthDB{gorm: db}
}

func (db *roleAuthDB) StoreRoleAuth(item *RoleAuth) error {
	if err := db.gorm.Create(item).Error; err != nil {
		log.Error("StoreRoleAuth error", "err", err)
		return err
	}
	return nil
}

func (db *roleAuthDB) StoreRoleAuths(list []*RoleAuth) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreRoleAuths error", "err", err)
		return err
	}
	return nil
}

func (db *roleAuthDB) GetByRoleID(roleID int64) ([]*RoleAuth, error) {
	var list []*RoleAuth
	if err := db.gorm.Where("role_id = ?", roleID).Find(&list).Error; err != nil {
		log.Error("GetByRoleID RoleAuth error", "err", err)
		return nil, err
	}
	return list, nil
}

func (db *roleAuthDB) GetByAuthID(authID int64) ([]*RoleAuth, error) {
	var list []*RoleAuth
	if err := db.gorm.Where("auth_id = ?", authID).Find(&list).Error; err != nil {
		log.Error("GetByAuthID RoleAuth error", "err", err)
		return nil, err
	}
	return list, nil
}

func (db *roleAuthDB) DeleteByRoleID(roleID int64) error {
	if roleID <= 0 {
		return fmt.Errorf("invalid roleID")
	}
	if err := db.gorm.Where("role_id = ?", roleID).Delete(&RoleAuth{}).Error; err != nil {
		log.Error("DeleteByRoleID error", "err", err)
		return err
	}
	return nil
}

func (db *roleAuthDB) DeleteByAuthID(authID int64) error {
	if authID <= 0 {
		return fmt.Errorf("invalid authID")
	}
	if err := db.gorm.Where("auth_id = ?", authID).Delete(&RoleAuth{}).Error; err != nil {
		log.Error("DeleteByAuthID error", "err", err)
		return err
	}
	return nil
}
