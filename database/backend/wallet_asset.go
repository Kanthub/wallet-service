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
	WalletUUID string    `gorm:"column:wallet_uuid;type:varchar(255);not null;index" json:"wallet_uuid"`
	TokenID    string    `gorm:"column:token_id;type:varchar(255);default:''" json:"token_id"`
	ChainID    string    `gorm:"column:chain_id;type:varchar(255);default:'';index" json:"chain_id"`
	Balance    string    `gorm:"column:balance;type:numeric(78,0);not null" json:"balance"` // 使用 string 存储大数字（支持 uint256）
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
	GetByTokenChain(tokenID, chainID string) (*WalletAsset, error)

	GetByWalletUUID(walletUUID string) ([]*WalletAsset, error)
	GetByWalletTokenChain(walletUUID, tokenID, chainID string) (*WalletAsset, error)
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

func (db *walletAssetDB) GetByTokenChain(tokenID, chainID string) (*WalletAsset, error) {
	var a WalletAsset
	if err := db.gorm.Where("token_id = ? AND chain_id = ?", tokenID, chainID).First(&a).Error; err != nil {
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

func (db *walletAssetDB) GetByWalletUUID(walletUUID string) ([]*WalletAsset, error) {
	var list []*WalletAsset
	if err := db.gorm.
		Where("wallet_uuid = ?", walletUUID).
		Find(&list).Error; err != nil {
		log.Error("GetByWalletUUID WalletAsset error", "err", err)
		return nil, err
	}
	return list, nil
}

func (db *walletAssetDB) GetByWalletTokenChain(
	walletUUID, tokenID, chainID string,
) (*WalletAsset, error) {

	var a WalletAsset
	if err := db.gorm.
		Where("wallet_uuid = ? AND token_id = ? AND chain_id = ?",
			walletUUID, tokenID, chainID).
		First(&a).Error; err != nil {

		log.Error("GetByWalletTokenChain error", "err", err)
		return nil, err
	}
	return &a, nil
}
