// chain_token.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type ChainToken struct {
	Guid       string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	ChainUUID  string    `gorm:"column:chain_uuid;type:varchar(255);default:''" json:"chain_uuid"`
	TokenUUID  string    `gorm:"column:token_uuid;type:varchar(255);not null" json:"token_uuid"`
	CreateTime time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (ChainToken) TableName() string {
	return "chain_token"
}

type ChainTokenView interface {
	GetByGuid(guid string) (*ChainToken, error)
	GetByChainUUID(chainUUID string) ([]*ChainToken, error)
	GetByTokenUUID(tokenUUID string) ([]*ChainToken, error)
}

type ChainTokenDB interface {
	ChainTokenView

	StoreChainToken(item *ChainToken) error
	StoreChainTokens(list []*ChainToken) error
	UpdateChainToken(guid string, updates map[string]interface{}) error
}

type chainTokenDB struct {
	gorm *gorm.DB
}

func NewChainTokenDB(db *gorm.DB) ChainTokenDB {
	return &chainTokenDB{gorm: db}
}

func (db *chainTokenDB) StoreChainToken(item *ChainToken) error {
	if err := db.gorm.Create(item).Error; err != nil {
		log.Error("StoreChainToken error", "err", err)
		return err
	}
	return nil
}

func (db *chainTokenDB) StoreChainTokens(list []*ChainToken) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreChainTokens error", "err", err)
		return err
	}
	return nil
}

func (db *chainTokenDB) GetByGuid(guid string) (*ChainToken, error) {
	var item ChainToken
	if err := db.gorm.Where("guid = ?", guid).First(&item).Error; err != nil {
		log.Error("GetByGuid ChainToken error", "err", err)
		return nil, err
	}
	return &item, nil
}

func (db *chainTokenDB) GetByChainUUID(chainUUID string) ([]*ChainToken, error) {
	var list []*ChainToken
	if err := db.gorm.Where("chain_uuid = ?", chainUUID).Find(&list).Error; err != nil {
		log.Error("GetByChainUUID ChainToken error", "err", err)
		return nil, err
	}
	return list, nil
}

func (db *chainTokenDB) GetByTokenUUID(tokenUUID string) ([]*ChainToken, error) {
	var list []*ChainToken
	if err := db.gorm.Where("token_uuid = ?", tokenUUID).Find(&list).Error; err != nil {
		log.Error("GetByTokenUUID ChainToken error", "err", err)
		return nil, err
	}
	return list, nil
}

func (db *chainTokenDB) UpdateChainToken(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&ChainToken{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateChainToken error", "err", err)
		return err
	}
	return nil
}
