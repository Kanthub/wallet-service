package backend

import (
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/log"
)

type Auth struct {
	ID         int    `gorm:"primaryKey;column:id" json:"id"`
	AuthName   string `gorm:"type:varchar(255);default:''" json:"auth_name"` // 权限名称
	AuthUrl    string `gorm:"type:varchar(255);default:''" json:"auth_url"`  // 权限路径/接口地址
	UserID     int    `gorm:"type:int;default:0" json:"user_id"`             // 所属用户/管理员ID
	Pid        int    `gorm:"type:int;default:0" json:"pid"`                 // 父级权限ID
	Sort       int    `gorm:"type:int;default:0" json:"sort"`                // 排序
	Icon       string `gorm:"type:varchar(255);default:''" json:"icon"`      // 图标
	IsShow     int    `gorm:"type:int;default:1" json:"is_show"`             // 是否显示(1显示;0隐藏)
	Status     int    `gorm:"type:int;default:1" json:"status"`              // 状态(1启用;0禁用)
	CreateID   int    `gorm:"type:int;default:0" json:"create_id"`           // 创建人ID
	UpdateID   int    `gorm:"type:int;default:0" json:"update_id"`           // 修改人ID
	CreateTime int64  `gorm:"type:bigint;default:0" json:"create_time"`      // 创建时间(Unix时间戳)
	UpdateTime int64  `gorm:"type:bigint;default:0" json:"update_time"`      // 更新时间(Unix时间戳)
}

func (Auth) TableName() string {
	return "auth"
}

type AuthView interface {
	GetByID(id int) (*Auth, error)
	GetByUserID(userID int) ([]*Auth, error)
	GetByPID(pid int) ([]*Auth, error)
}

type AuthDB interface {
	AuthView

	StoreAuth(auth *Auth) error
	StoreAuths(auths []*Auth) error
}

type authDB struct {
	gorm *gorm.DB
}

func NewAuthDB(db *gorm.DB) AuthDB {
	return &authDB{gorm: db}
}

func (db *authDB) StoreAuth(auth *Auth) error {
	if err := db.gorm.Create(auth).Error; err != nil {
		log.Error("StoreAuth error:", err)
		return err
	}
	return nil
}

func (db *authDB) StoreAuths(auths []*Auth) error {
	if err := db.gorm.CreateInBatches(auths, len(auths)).Error; err != nil {
		log.Error("StoreAuths error:", err)
		return err
	}
	return nil
}

func (db *authDB) GetByID(id int) (*Auth, error) {
	var auth Auth
	if err := db.gorm.First(&auth, id).Error; err != nil {
		log.Error("GetByID error:", err)
		return nil, err
	}
	return &auth, nil
}

func (db *authDB) GetByUserID(userID int) ([]*Auth, error) {
	var list []*Auth
	if err := db.gorm.Where("user_id = ?", userID).Find(&list).Error; err != nil {
		log.Error("GetByUserID error:", err)
		return nil, err
	}
	return list, nil
}

func (db *authDB) GetByPID(pid int) ([]*Auth, error) {
	var list []*Auth
	if err := db.gorm.Where("pid = ?", pid).Find(&list).Error; err != nil {
		log.Error("GetByPID error:", err)
		return nil, err
	}
	return list, nil
}
