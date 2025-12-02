// fiat_currency_rate.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type FiatCurrencyRate struct {
	Guid       string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	KeyName    string    `gorm:"column:key_name;type:varchar(255);not null" json:"key_name"`
	ValueData  string    `gorm:"column:value_data;type:varchar(255);not null" json:"value_data"`
	CreateTime time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (FiatCurrencyRate) TableName() string {
	return "fiat_currency_rate"
}

type FiatCurrencyRateView interface {
	GetByGuid(guid string) (*FiatCurrencyRate, error)
	GetByKeyName(key string) (*FiatCurrencyRate, error)
}

type FiatCurrencyRateDB interface {
	FiatCurrencyRateView

	StoreFiatCurrencyRate(r *FiatCurrencyRate) error
	StoreFiatCurrencyRates(list []*FiatCurrencyRate) error
	UpdateFiatCurrencyRate(guid string, updates map[string]interface{}) error
}

type fiatCurrencyRateDB struct {
	gorm *gorm.DB
}

func NewFiatCurrencyRateDB(db *gorm.DB) FiatCurrencyRateDB {
	return &fiatCurrencyRateDB{gorm: db}
}

func (db *fiatCurrencyRateDB) StoreFiatCurrencyRate(r *FiatCurrencyRate) error {
	if err := db.gorm.Create(r).Error; err != nil {
		log.Error("StoreFiatCurrencyRate error", "err", err)
		return err
	}
	return nil
}

func (db *fiatCurrencyRateDB) StoreFiatCurrencyRates(list []*FiatCurrencyRate) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreFiatCurrencyRates error", "err", err)
		return err
	}
	return nil
}

func (db *fiatCurrencyRateDB) GetByGuid(guid string) (*FiatCurrencyRate, error) {
	var r FiatCurrencyRate
	if err := db.gorm.Where("guid = ?", guid).First(&r).Error; err != nil {
		log.Error("GetByGuid FiatCurrencyRate error", "err", err)
		return nil, err
	}
	return &r, nil
}

func (db *fiatCurrencyRateDB) GetByKeyName(key string) (*FiatCurrencyRate, error) {
	var r FiatCurrencyRate
	if err := db.gorm.Where("key_name = ?", key).First(&r).Error; err != nil {
		log.Error("GetByKeyName FiatCurrencyRate error", "err", err)
		return nil, err
	}
	return &r, nil
}

func (db *fiatCurrencyRateDB) UpdateFiatCurrencyRate(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&FiatCurrencyRate{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateFiatCurrencyRate error", "err", err)
		return err
	}
	return nil
}
