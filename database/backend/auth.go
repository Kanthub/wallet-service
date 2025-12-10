package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type Auth struct {
	Guid       string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	AuthName   string    `gorm:"column:auth_name;type:varchar(255);default:''" json:"auth_name"`
	AuthURL    string    `gorm:"column:auth_url;type:varchar(255);default:''" json:"auth_url"`
	UserID     int64     `gorm:"column:user_id;type:int;default:0" json:"user_id"`
	PID        int64     `gorm:"column:pid;type:int;default:0" json:"pid"`
	Sort       int64     `gorm:"column:sort;type:int;default:0" json:"sort"`
	Icon       string    `gorm:"column:icon;type:varchar(255);default:''" json:"icon"`
	IsShow     int64     `gorm:"column:is_show;type:int;default:1" json:"is_show"`
	Status     int64     `gorm:"column:status;type:int;default:1" json:"status"`
	CreateID   int64     `gorm:"column:create_id;type:int;default:0" json:"create_id"`
	UpdateID   int64     `gorm:"column:update_id;type:int;default:0" json:"update_id"`
	CreateTime time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (Auth) TableName() string {
	return "auth"
}

type AuthView interface {
	GetByGuid(guid string) (*Auth, error)
	GetAuthList(page, pageSize int, filters map[string]interface{}) ([]*Auth, int64, error)
}

type AuthDB interface {
	AuthView

	StoreAuth(auth *Auth) error
	StoreAuths(list []*Auth) error
	UpdateAuth(guid string, updates map[string]interface{}) error
}

type authDB struct {
	gorm *gorm.DB
}

func NewAuthDB(db *gorm.DB) AuthDB {
	return &authDB{gorm: db}
}

func (db *authDB) StoreAuth(a *Auth) error {
	if err := db.gorm.Create(a).Error; err != nil {
		log.Error("StoreAuth error", "err", err)
		return err
	}
	return nil
}

func (db *authDB) StoreAuths(list []*Auth) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreAuths error", "err", err)
		return err
	}
	return nil
}

func (db *authDB) GetByGuid(guid string) (*Auth, error) {
	var a Auth
	if err := db.gorm.Where("guid = ?", guid).First(&a).Error; err != nil {
		log.Error("GetByGuid auth error", "err", err)
		return nil, err
	}
	return &a, nil
}

func (db *authDB) GetAuthList(page, pageSize int, filters map[string]interface{}) ([]*Auth, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var list []*Auth
	query := db.gorm.Model(&Auth{})

	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}
		switch key {
		case "auth_name", "auth_url":
			query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error("GetAuthList count error", "err", err)
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
		log.Error("GetAuthList list error", "err", err)
		return nil, 0, err
	}

	return list, total, nil
}

func (db *authDB) UpdateAuth(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}
	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&Auth{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateAuth error", "err", err)
		return err
	}
	return nil
}
