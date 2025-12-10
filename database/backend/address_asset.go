// address_asset.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type AddressAsset struct {
	Guid        string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	TokenID     string    `gorm:"column:token_id;type:varchar(255);default:''" json:"token_id"`
	WalletUUID  string    `gorm:"column:wallet_uuid;type:varchar(255);default:''" json:"wallet_uuid"`
	AddressUUID string    `gorm:"column:address_uuid;type:varchar(255);default:''" json:"address_uuid"`
	AssetUsdt   string    `gorm:"column:asset_usdt;type:numeric(20,8);not null" json:"asset_usdt"`
	AssetUsd    string    `gorm:"column:asset_usd;type:numeric(20,8);not null" json:"asset_usd"`
	Balance     string    `gorm:"column:balance;type:numeric(78,0);not null" json:"balance"` // 使用 string 存储大数字（支持 uint256）
	CreateTime  time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (AddressAsset) TableName() string {
	return "address_asset"
}

type AddressAssetView interface {
	GetByGuid(guid string) (*AddressAsset, error)
	GetByAddressUUID(addressUUID string) ([]*AddressAsset, error)
}

type AddressAssetDB interface {
	AddressAssetView

	StoreAddressAsset(a *AddressAsset) error
	StoreAddressAssets(list []*AddressAsset) error
	UpdateAddressAsset(guid string, updates map[string]interface{}) error
}

type addressAssetDB struct {
	gorm *gorm.DB
}

func NewAddressAssetDB(db *gorm.DB) AddressAssetDB {
	return &addressAssetDB{gorm: db}
}

func (db *addressAssetDB) StoreAddressAsset(a *AddressAsset) error {
	if err := db.gorm.Create(a).Error; err != nil {
		log.Error("StoreAddressAsset error", "err", err)
		return err
	}
	return nil
}

func (db *addressAssetDB) StoreAddressAssets(list []*AddressAsset) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreAddressAssets error", "err", err)
		return err
	}
	return nil
}

func (db *addressAssetDB) GetByGuid(guid string) (*AddressAsset, error) {
	var a AddressAsset
	if err := db.gorm.Where("guid = ?", guid).First(&a).Error; err != nil {
		log.Error("GetByGuid AddressAsset error", "err", err)
		return nil, err
	}
	return &a, nil
}

func (db *addressAssetDB) GetByAddressUUID(addressUUID string) ([]*AddressAsset, error) {
	var list []*AddressAsset
	if err := db.gorm.Where("address_uuid = ?", addressUUID).Find(&list).Error; err != nil {
		log.Error("GetByAddressUUID AddressAsset error", "err", err)
		return nil, err
	}
	return list, nil
}

func (db *addressAssetDB) UpdateAddressAsset(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&AddressAsset{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateAddressAsset error", "err", err)
		return err
	}
	return nil
}
