// wallet_address.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type WalletAddress struct {
	Guid         string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	AddressIndex int64     `gorm:"column:address_index;type:integer;check:address_index > 0" json:"address_index"`
	Address      string    `gorm:"column:address;type:varchar(70);not null;index" json:"address"`
	WalletUUID   string    `gorm:"column:wallet_uuid;type:varchar(255);default:'';index" json:"wallet_uuid"`
	ChainID      string    `gorm:"column:chain_id;type:varchar(255);default:'';index" json:"chain_id"`
	CreateTime   time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (WalletAddress) TableName() string {
	return "wallet_address"
}

type WalletAddressView interface {
	GetByGuid(guid string) (*WalletAddress, error)
	GetByAddress(address string) (*WalletAddress, error)
	GetByWalletUUID(walletUUID string) ([]*WalletAddress, error)
}

type WalletAddressDB interface {
	WalletAddressView

	StoreWalletAddress(a *WalletAddress) error
	StoreWalletAddresses(list []*WalletAddress) error
	UpdateWalletAddress(guid string, updates map[string]interface{}) error
}

type walletAddressDB struct {
	gorm *gorm.DB
}

func NewWalletAddressDB(db *gorm.DB) WalletAddressDB {
	return &walletAddressDB{gorm: db}
}

func (db *walletAddressDB) StoreWalletAddress(a *WalletAddress) error {
	if err := db.gorm.Create(a).Error; err != nil {
		log.Error("StoreWalletAddress error", "err", err)
		return err
	}
	return nil
}

func (db *walletAddressDB) StoreWalletAddresses(list []*WalletAddress) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreWalletAddresses error", "err", err)
		return err
	}
	return nil
}

func (db *walletAddressDB) GetByGuid(guid string) (*WalletAddress, error) {
	var a WalletAddress
	if err := db.gorm.Where("guid = ?", guid).First(&a).Error; err != nil {
		log.Error("GetByGuid WalletAddress error", "err", err)
		return nil, err
	}
	return &a, nil
}

func (db *walletAddressDB) GetByAddress(address string) (*WalletAddress, error) {
	var a WalletAddress
	if err := db.gorm.Where("address = ?", address).First(&a).Error; err != nil {
		log.Error("GetByAddress WalletAddress error", "err", err)
		return nil, err
	}
	return &a, nil
}

func (db *walletAddressDB) GetByWalletUUID(walletUUID string) ([]*WalletAddress, error) {
	var list []*WalletAddress
	if err := db.gorm.Where("wallet_uuid = ?", walletUUID).Find(&list).Error; err != nil {
		log.Error("GetByWalletUUID WalletAddress error", "err", err)
		return nil, err
	}
	return list, nil
}

func (db *walletAddressDB) UpdateWalletAddress(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&WalletAddress{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateWalletAddress error", "err", err)
		return err
	}
	return nil
}
