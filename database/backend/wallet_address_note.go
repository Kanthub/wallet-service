// wallet_address_note.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type WalletAddressNote struct {
	Guid       string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	DeviceUUID string    `gorm:"column:device_uuid;type:varchar(255);not null" json:"device_uuid"`
	ChainUUID  string    `gorm:"column:chain_uuid;type:varchar(255);default:''" json:"chain_uuid"`
	Memo       string    `gorm:"column:memo;type:varchar(255);not null" json:"memo"`
	Address    string    `gorm:"column:address;type:varchar(255);not null" json:"address"`
	CreateTime time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (WalletAddressNote) TableName() string {
	return "wallet_address_note"
}

type WalletAddressNoteView interface {
	GetByGuid(guid string) (*WalletAddressNote, error)
	GetByDeviceUUID(deviceUUID string) ([]*WalletAddressNote, error)
}

type WalletAddressNoteDB interface {
	WalletAddressNoteView

	StoreWalletAddressNote(n *WalletAddressNote) error
	StoreWalletAddressNotes(list []*WalletAddressNote) error
	UpdateWalletAddressNote(guid string, updates map[string]interface{}) error
}

type walletAddressNoteDB struct {
	gorm *gorm.DB
}

func NewWalletAddressNoteDB(db *gorm.DB) WalletAddressNoteDB {
	return &walletAddressNoteDB{gorm: db}
}

func (db *walletAddressNoteDB) StoreWalletAddressNote(n *WalletAddressNote) error {
	if err := db.gorm.Create(n).Error; err != nil {
		log.Error("StoreWalletAddressNote error", "err", err)
		return err
	}
	return nil
}

func (db *walletAddressNoteDB) StoreWalletAddressNotes(list []*WalletAddressNote) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreWalletAddressNotes error", "err", err)
		return err
	}
	return nil
}

func (db *walletAddressNoteDB) GetByGuid(guid string) (*WalletAddressNote, error) {
	var n WalletAddressNote
	if err := db.gorm.Where("guid = ?", guid).First(&n).Error; err != nil {
		log.Error("GetByGuid WalletAddressNote error", "err", err)
		return nil, err
	}
	return &n, nil
}

func (db *walletAddressNoteDB) GetByDeviceUUID(deviceUUID string) ([]*WalletAddressNote, error) {
	var list []*WalletAddressNote
	if err := db.gorm.Where("device_uuid = ?", deviceUUID).Find(&list).Error; err != nil {
		log.Error("GetByDeviceUUID WalletAddressNote error", "err", err)
		return nil, err
	}
	return list, nil
}

func (db *walletAddressNoteDB) UpdateWalletAddressNote(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&WalletAddressNote{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateWalletAddressNote error", "err", err)
		return err
	}
	return nil
}
