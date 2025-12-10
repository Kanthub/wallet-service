// kline.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type Kline struct {
	Guid         string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	TokenID      string    `gorm:"column:token_id;type:varchar;not null" json:"token_id"`
	TimeInterval string    `gorm:"column:time_interval;type:varchar;not null" json:"time_interval"`
	OpenTime     time.Time `gorm:"column:open_time;type:timestamp;not null" json:"open_time"`
	OpenPrice    string    `gorm:"column:open_price;type:numeric(20,8);not null" json:"open_price"`
	HighPrice    string    `gorm:"column:high_price;type:numeric(20,8);not null" json:"high_price"`
	LowPrice     string    `gorm:"column:low_price;type:numeric(20,8);not null" json:"low_price"`
	ClosePrice   string    `gorm:"column:close_price;type:numeric(20,8);not null" json:"close_price"`
	Volume       string    `gorm:"column:volume;type:numeric;not null" json:"volume"` // UINT256
	QuoteVolume  string    `gorm:"column:quote_volume;type:numeric;default:0" json:"quote_volume"`
	TradeCount   string    `gorm:"column:trade_count;type:numeric;default:0" json:"trade_count"`
	CreateTime   time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (Kline) TableName() string {
	return "kline"
}

type KlineView interface {
	GetByGuid(guid string) (*Kline, error)
	GetByKey(tokenID, interval string, openTime time.Time) (*Kline, error)
}

type KlineDB interface {
	KlineView

	StoreKline(k *Kline) error
	StoreKlines(list []*Kline) error
	UpdateKline(guid string, updates map[string]interface{}) error
}

type klineDB struct {
	gorm *gorm.DB
}

func NewKlineDB(db *gorm.DB) KlineDB {
	return &klineDB{gorm: db}
}

func (db *klineDB) StoreKline(k *Kline) error {
	if err := db.gorm.Create(k).Error; err != nil {
		log.Error("StoreKline error", "err", err)
		return err
	}
	return nil
}

func (db *klineDB) StoreKlines(list []*Kline) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreKlines error", "err", err)
		return err
	}
	return nil
}

func (db *klineDB) GetByGuid(guid string) (*Kline, error) {
	var k Kline
	if err := db.gorm.Where("guid = ?", guid).First(&k).Error; err != nil {
		log.Error("GetByGuid Kline error", "err", err)
		return nil, err
	}
	return &k, nil
}

func (db *klineDB) GetByKey(tokenID, interval string, openTime time.Time) (*Kline, error) {
	var k Kline
	if err := db.gorm.Where("token_id = ? AND time_interval = ? AND open_time = ?", tokenID, interval, openTime).First(&k).Error; err != nil {
		log.Error("GetByKey Kline error", "err", err)
		return nil, err
	}
	return &k, nil
}

func (db *klineDB) UpdateKline(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&Kline{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateKline error", "err", err)
		return err
	}
	return nil
}
