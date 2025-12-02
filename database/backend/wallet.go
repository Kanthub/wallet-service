// wallet.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type Wallet struct {
	Guid       string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	DeviceUUID string    `gorm:"column:device_uuid;type:varchar(255);not null" json:"device_uuid"`
	WalletUUID string    `gorm:"column:wallet_uuid;type:varchar(255);not null" json:"wallet_uuid"`
	ChainUUID  string    `gorm:"column:chain_uuid;type:varchar(255);default:''" json:"chain_uuid"`
	WalletName string    `gorm:"column:wallet_name;type:varchar(70);default:'roothash'" json:"wallet_name"`
	AssetUsdt  string    `gorm:"column:asset_usdt;type:numeric(20,8);not null" json:"asset_usdt"`
	AssetUsd   string    `gorm:"column:asset_usd;type:numeric(20,8);not null" json:"asset_usd"`
	CreateTime time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (Wallet) TableName() string {
	return "wallet"
}

type WalletView interface {
	GetByGuid(guid string) (*Wallet, error)
	GetByWalletUUID(walletUUID string) (*Wallet, error)
	GetWalletList(page, pageSize int, filters map[string]interface{}) ([]*Wallet, int64, error)
}

type WalletDB interface {
	WalletView

	StoreWallet(w *Wallet) error
	StoreWallets(list []*Wallet) error
	UpdateWallet(guid string, updates map[string]interface{}) error
}

type walletDB struct {
	gorm *gorm.DB
}

func NewWalletDB(db *gorm.DB) WalletDB {
	return &walletDB{gorm: db}
}

func (db *walletDB) StoreWallet(w *Wallet) error {
	if err := db.gorm.Create(w).Error; err != nil {
		log.Error("StoreWallet error", "err", err)
		return err
	}
	return nil
}

func (db *walletDB) StoreWallets(list []*Wallet) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreWallets error", "err", err)
		return err
	}
	return nil
}

func (db *walletDB) GetByGuid(guid string) (*Wallet, error) {
	var w Wallet
	if err := db.gorm.Where("guid = ?", guid).First(&w).Error; err != nil {
		log.Error("GetByGuid wallet error", "err", err)
		return nil, err
	}
	return &w, nil
}

func (db *walletDB) GetByWalletUUID(walletUUID string) (*Wallet, error) {
	var w Wallet
	if err := db.gorm.Where("wallet_uuid = ?", walletUUID).First(&w).Error; err != nil {
		log.Error("GetByWalletUUID wallet error", "err", err)
		return nil, err
	}
	return &w, nil
}

func (db *walletDB) GetWalletList(page, pageSize int, filters map[string]interface{}) ([]*Wallet, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	var list []*Wallet
	query := db.gorm.Model(&Wallet{})

	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}
		switch key {
		case "wallet_name":
			query = query.Where("wallet_name LIKE ?", "%"+value.(string)+"%")
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error("GetWalletList count error", "err", err)
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
		log.Error("GetWalletList list error", "err", err)
		return nil, 0, err
	}

	return list, total, nil
}

func (db *walletDB) UpdateWallet(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&Wallet{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateWallet error", "err", err)
		return err
	}
	return nil
}
