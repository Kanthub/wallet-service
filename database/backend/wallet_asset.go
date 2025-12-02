// wallet_asset.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type WalletAsset struct {
	Guid       string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	TokenUUID  string    `gorm:"column:token_uuid;type:varchar(255);default:''" json:"token_uuid"`
	ChainUUID  string    `gorm:"column:chain_uuid;type:varchar(255);default:''" json:"chain_uuid"`
	Balance    int64     `gorm:"column:balance;type:integer" json:"balance"`
	AssetUsdt  string    `gorm:"column:asset_usdt;type:numeric(20,8);not null" json:"asset_usdt"`
	AssetUsd   string    `gorm:"column:asset_usd;type:numeric(20,8);not null" json:"asset_usd"`
	CreateTime time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (WalletAsset) TableName() string {
	return "wallet_asset"
}

type WalletAssetView interface {
	GetByGuid(guid string) (*WalletAsset, error)
	GetByTokenChain(tokenUUID, chainUUID string) (*WalletAsset, error)
}

type WalletAssetDB interface {
	WalletAssetView

	StoreWalletAsset(a *WalletAsset) error
	StoreWalletAssets(list []*WalletAsset) error
	UpdateWalletAsset(guid string, updates map[string]interface{}) error
}

type walletAssetDB struct {
	gorm *gorm.DB
}

func NewWalletAssetDB(db *gorm.DB) WalletAssetDB {
	return &walletAssetDB{gorm: db}
}

func (db *walletAssetDB) StoreWalletAsset(a *WalletAsset) error {
	if err := db.gorm.Create(a).Error; err != nil {
		log.Error("StoreWalletAsset error", "err", err)
		return err
	}
	return nil
}

func (db *walletAssetDB) StoreWalletAssets(list []*WalletAsset) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreWalletAssets error", "err", err)
		return err
	}
	return nil
}

func (db *walletAssetDB) GetByGuid(guid string) (*WalletAsset, error) {
	var a WalletAsset
	if err := db.gorm.Where("guid = ?", guid).First(&a).Error; err != nil {
		log.Error("GetByGuid WalletAsset error", "err", err)
		return nil, err
	}
	return &a, nil
}

func (db *walletAssetDB) GetByTokenChain(tokenUUID, chainUUID string) (*WalletAsset, error) {
	var a WalletAsset
	if err := db.gorm.Where("token_uuid = ? AND chain_uuid = ?", tokenUUID, chainUUID).First(&a).Error; err != nil {
		log.Error("GetByTokenChain WalletAsset error", "err", err)
		return nil, err
	}
	return &a, nil
}

func (db *walletAssetDB) UpdateWalletAsset(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&WalletAsset{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateWalletAsset error", "err", err)
		return err
	}
	return nil
}
