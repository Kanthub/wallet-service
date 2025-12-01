package backend

import (
	"fmt"
	"gorm.io/gorm"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

type Admin struct {
	ID         int64     `gorm:"primaryKey;column:id" json:"id"`
	LoginName  string    `gorm:"type:varchar(32);not null;unique" json:"login_name"` // 登录名
	RealName   string    `gorm:"type:varchar(32);unique" json:"real_name"`           // 真实姓名
	Password   string    `gorm:"type:varchar(100);not null" json:"password"`         // 密码(加密后)
	RoleIDs    string    `gorm:"type:varchar(255);default:''" json:"role_ids"`       // 角色 ID 列表（JSON/CSV）
	Phone      string    `gorm:"type:varchar(11);unique" json:"phone"`               // 手机号
	Email      string    `gorm:"type:varchar(32)" json:"email"`                      // 邮箱
	Salt       string    `gorm:"type:varchar(255);default:''" json:"salt"`           // 密码盐
	LastLogin  int64     `gorm:"type:bigint;default:0" json:"last_login"`            // 最后登录时间戳
	LastIP     string    `gorm:"type:varchar(255);default:''" json:"last_ip"`        // 最后登录 IP
	Status     int       `gorm:"type:int;default:1" json:"status"`                   // 状态(1启用;0禁用)
	CreateID   int       `gorm:"type:int;default:0" json:"create_id"`                // 创建人
	UpdateID   int       `gorm:"type:int;default:0" json:"update_id"`                // 修改人
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`                  // 创建时间
	UpdateTime time.Time `gorm:"autoUpdateTime" json:"update_time"`                  // 更新时间
}

func (Admin) TableName() string {
	return "admin"
}

type AdminView interface {
	GetByID(id int64) (*Admin, error)
	GetByStatus(status int) ([]*Admin, error)
	GetByLoginName(loginName string) (*Admin, error)
	GetAdminList(page, pageSize int, filters map[string]interface{}) ([]*Admin, int64, error)
}

type AdminDB interface {
	AdminView

	StoreAdmin(admin *Admin) error
	StoreAdmins(admins []*Admin) error
	UpdateAdmin(id int64, updates map[string]interface{}) error
}

type adminDB struct {
	gorm *gorm.DB
}

func NewAdminDB(db *gorm.DB) AdminDB {
	return &adminDB{gorm: db}
}

func (db *adminDB) StoreAdmin(admin *Admin) error {
	if err := db.gorm.Create(admin).Error; err != nil {
		log.Error("store admin error:", err)
		return err
	}
	return nil
}

func (db *adminDB) StoreAdmins(admins []*Admin) error {
	if err := db.gorm.CreateInBatches(admins, len(admins)).Error; err != nil {
		log.Error("store admins error:", err)
		return err
	}
	return nil
}

func (db *adminDB) GetByID(id int64) (*Admin, error) {
	var admin Admin
	if err := db.gorm.First(&admin, id).Error; err != nil {
		log.Error("GetByID error:", err)
		return nil, err
	}
	return &admin, nil
}

func (db *adminDB) GetByStatus(status int) ([]*Admin, error) {
	var list []*Admin
	if err := db.gorm.Where("status = ?", status).Find(&list).Error; err != nil {
		log.Error("get by status error:", err)
		return nil, err
	}
	return list, nil
}

func (db *adminDB) GetByLoginName(loginName string) (*Admin, error) {
	var admin Admin
	if err := db.gorm.Where("login_name = ?", loginName).First(&admin).Error; err != nil {
		log.Error("get login name error:", err)
		return nil, err
	}
	return &admin, nil
}

func (db *adminDB) GetAdminList(page, pageSize int, filters map[string]interface{}) ([]*Admin, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var list []*Admin
	query := db.gorm.Model(&Admin{})

	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}

		switch key {
		case "login_name":
			query = query.Where("login_name LIKE ?", "%"+value.(string)+"%")
		case "real_name":
			query = query.Where("real_name LIKE ?", "%"+value.(string)+"%")
		case "email":
			query = query.Where("email LIKE ?", "%"+value.(string)+"%")
		case "phone":
			query = query.Where("phone LIKE ?", "%"+value.(string)+"%")
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error("GetAdminList count error:", err)
		return nil, 0, err
	}

	if err := query.Order("id DESC").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
		log.Error("get admin list error:", err)
		return nil, 0, err
	}

	return list, total, nil
}

func (db *adminDB) UpdateAdmin(id int64, updates map[string]interface{}) error {
	if id <= 0 {
		return fmt.Errorf("invalid id")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["update_time"] = time.Now()

	if err := db.gorm.Model(&Admin{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		log.Error("UpdateAdmin error:", err)
		return err
	}
	return nil
}
