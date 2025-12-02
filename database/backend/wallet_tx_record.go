// wallet_tx_record.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type WalletTxRecord struct {
	Guid        string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	TxTime      string    `gorm:"column:tx_time;type:varchar(500);not null" json:"tx_time"`
	ChainUUID   string    `gorm:"column:chain_uuid;type:varchar(255);default:''" json:"chain_uuid"`
	TokenUUID   string    `gorm:"column:token_uuid;type:varchar(255);default:''" json:"token_uuid"`
	FromAddress string    `gorm:"column:from_address;type:varchar(70);not null" json:"from_address"`
	ToAddress   string    `gorm:"column:to_address;type:varchar(70);not null" json:"to_address"`
	Amount      int64     `gorm:"column:amount;type:integer" json:"amount"`
	Memo        string    `gorm:"column:memo;type:varchar(500);not null" json:"memo"`
	Hash        string    `gorm:"column:hash;type:varchar(500);not null" json:"hash"`
	BlockHeight string    `gorm:"column:block_height;type:varchar(500);not null" json:"block_height"`
	ExplorerURL string    `gorm:"column:explorer_url;type:varchar(500);not null" json:"explorer_url"`
	CreateTime  time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (WalletTxRecord) TableName() string {
	return "wallet_tx_record"
}

type WalletTxRecordView interface {
	GetByGuid(guid string) (*WalletTxRecord, error)
	GetByHash(hash string) (*WalletTxRecord, error)
	GetTxList(page, pageSize int, filters map[string]interface{}) ([]*WalletTxRecord, int64, error)
}

type WalletTxRecordDB interface {
	WalletTxRecordView

	StoreWalletTxRecord(r *WalletTxRecord) error
	StoreWalletTxRecords(list []*WalletTxRecord) error
	UpdateWalletTxRecord(guid string, updates map[string]interface{}) error
}

type walletTxRecordDB struct {
	gorm *gorm.DB
}

func NewWalletTxRecordDB(db *gorm.DB) WalletTxRecordDB {
	return &walletTxRecordDB{gorm: db}
}

func (db *walletTxRecordDB) StoreWalletTxRecord(r *WalletTxRecord) error {
	if err := db.gorm.Create(r).Error; err != nil {
		log.Error("StoreWalletTxRecord error", "err", err)
		return err
	}
	return nil
}

func (db *walletTxRecordDB) StoreWalletTxRecords(list []*WalletTxRecord) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreWalletTxRecords error", "err", err)
		return err
	}
	return nil
}

func (db *walletTxRecordDB) GetByGuid(guid string) (*WalletTxRecord, error) {
	var r WalletTxRecord
	if err := db.gorm.Where("guid = ?", guid).First(&r).Error; err != nil {
		log.Error("GetByGuid WalletTxRecord error", "err", err)
		return nil, err
	}
	return &r, nil
}

func (db *walletTxRecordDB) GetByHash(hash string) (*WalletTxRecord, error) {
	var r WalletTxRecord
	if err := db.gorm.Where("hash = ?", hash).First(&r).Error; err != nil {
		log.Error("GetByHash WalletTxRecord error", "err", err)
		return nil, err
	}
	return &r, nil
}

func (db *walletTxRecordDB) GetTxList(page, pageSize int, filters map[string]interface{}) ([]*WalletTxRecord, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var list []*WalletTxRecord
	query := db.gorm.Model(&WalletTxRecord{})

	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}
		switch key {
		case "from_address", "to_address", "hash":
			query = query.Where(key+" = ?", value)
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error("GetTxList count error", "err", err)
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
		log.Error("GetTxList list error", "err", err)
		return nil, 0, err
	}

	return list, total, nil
}

func (db *walletTxRecordDB) UpdateWalletTxRecord(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&WalletTxRecord{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateWalletTxRecord error", "err", err)
		return err
	}
	return nil
}
