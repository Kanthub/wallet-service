// market_price.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type MarketPrice struct {
	Guid        string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	ChainUUID   string    `gorm:"column:chain_uuid;type:varchar(255);default:''" json:"chain_uuid"`
	TokenUUID   string    `gorm:"column:token_uuid;type:varchar(255);default:''" json:"token_uuid"`
	UsdtPrice   string    `gorm:"column:usdt_price;type:numeric(20,8);not null" json:"usdt_price"`
	UsdPrice    string    `gorm:"column:usd_price;type:numeric(20,8);not null" json:"usd_price"`
	MarketCap   int64     `gorm:"column:market_cap;type:integer" json:"market_cap"`
	Liquidity   int64     `gorm:"column:liquidity;type:integer" json:"liquidity"`
	Volume24h   int64     `gorm:"column:24h_volume;type:integer" json:"24h_volume"`
	PriceChange string    `gorm:"column:price_change;type:varchar(255);not null" json:"price_change"`
	Ranking     string    `gorm:"column:ranking;type:varchar(255);not null" json:"ranking"`
	CreateTime  time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (MarketPrice) TableName() string {
	return "market_price"
}

type MarketPriceView interface {
	GetByGuid(guid string) (*MarketPrice, error)
	GetByTokenUUID(tokenUUID string) (*MarketPrice, error)
}

type MarketPriceDB interface {
	MarketPriceView

	StoreMarketPrice(m *MarketPrice) error
	StoreMarketPrices(list []*MarketPrice) error
	UpdateMarketPrice(guid string, updates map[string]interface{}) error
}

type marketPriceDB struct {
	gorm *gorm.DB
}

func NewMarketPriceDB(db *gorm.DB) MarketPriceDB {
	return &marketPriceDB{gorm: db}
}

func (db *marketPriceDB) StoreMarketPrice(m *MarketPrice) error {
	if err := db.gorm.Create(m).Error; err != nil {
		log.Error("StoreMarketPrice error", "err", err)
		return err
	}
	return nil
}

func (db *marketPriceDB) StoreMarketPrices(list []*MarketPrice) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreMarketPrices error", "err", err)
		return err
	}
	return nil
}

func (db *marketPriceDB) GetByGuid(guid string) (*MarketPrice, error) {
	var m MarketPrice
	if err := db.gorm.Where("guid = ?", guid).First(&m).Error; err != nil {
		log.Error("GetByGuid MarketPrice error", "err", err)
		return nil, err
	}
	return &m, nil
}

func (db *marketPriceDB) GetByTokenUUID(tokenUUID string) (*MarketPrice, error) {
	var m MarketPrice
	if err := db.gorm.Where("token_uuid = ?", tokenUUID).First(&m).Error; err != nil {
		log.Error("GetByTokenUUID MarketPrice error", "err", err)
		return nil, err
	}
	return &m, nil
}

func (db *marketPriceDB) UpdateMarketPrice(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&MarketPrice{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateMarketPrice error", "err", err)
		return err
	}
	return nil
}
