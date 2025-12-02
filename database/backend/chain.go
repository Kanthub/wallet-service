package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type Chain struct {
	Guid            string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	ChainName       string    `gorm:"column:chain_name;type:varchar(70);not null" json:"chain_name"`
	ChainMark       string    `gorm:"column:chain_mark;type:varchar(70);not null" json:"chain_mark"`
	ChainLogo       string    `gorm:"column:chain_logo;type:varchar(200);not null" json:"chain_logo"`
	ChainActiveLogo string    `gorm:"column:chain_active_logo;type:varchar(200);not null" json:"chain_active_logo"`
	ChainModelType  string    `gorm:"column:chain_model_type;type:varchar(10);not null" json:"chain_model_type"`
	CreateTime      time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (Chain) TableName() string {
	return "chain"
}

type ChainView interface {
	GetByGuid(guid string) (*Chain, error)
	GetByName(name string) (*Chain, error)
	GetChainList(page, pageSize int, filters map[string]interface{}) ([]*Chain, int64, error)
}

type ChainDB interface {
	ChainView

	StoreChain(chain *Chain) error
	StoreChains(list []*Chain) error
	UpdateChain(guid string, updates map[string]interface{}) error
}

type chainDB struct {
	gorm *gorm.DB
}

func NewChainDB(db *gorm.DB) ChainDB {
	return &chainDB{gorm: db}
}

func (db *chainDB) StoreChain(c *Chain) error {
	if err := db.gorm.Create(c).Error; err != nil {
		log.Error("StoreChain error", "err", err)
		return err
	}
	return nil
}

func (db *chainDB) StoreChains(list []*Chain) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreChains error", "err", err)
		return err
	}
	return nil
}

func (db *chainDB) GetByGuid(guid string) (*Chain, error) {
	var c Chain
	if err := db.gorm.Where("guid = ?", guid).First(&c).Error; err != nil {
		log.Error("GetByGuid chain error", "err", err)
		return nil, err
	}
	return &c, nil
}

func (db *chainDB) GetByName(name string) (*Chain, error) {
	var c Chain
	if err := db.gorm.Where("chain_name = ?", name).First(&c).Error; err != nil {
		log.Error("GetByName chain error", "err", err)
		return nil, err
	}
	return &c, nil
}

func (db *chainDB) GetChainList(page, pageSize int, filters map[string]interface{}) ([]*Chain, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var list []*Chain
	query := db.gorm.Model(&Chain{})

	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}
		switch key {
		case "chain_name", "chain_mark":
			query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error("GetChainList count error", "err", err)
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
		log.Error("GetChainList list error", "err", err)
		return nil, 0, err
	}

	return list, total, nil
}

func (db *chainDB) UpdateChain(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}
	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&Chain{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateChain error", "err", err)
		return err
	}
	return nil
}
