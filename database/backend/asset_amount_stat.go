// asset_amount_stat.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type AssetAmountStat struct {
	Guid       string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	AssetUUID  string    `gorm:"column:asset_uuid;type:varchar(255);default:''" json:"asset_uuid"`
	TimeDate   string    `gorm:"column:time_date;type:varchar(255);not null" json:"time_date"`
	Amount     string    `gorm:"column:amount;type:numeric(78,0);not null" json:"amount"` // 使用 string 存储大数字（支持 uint256）
	CreateTime time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (AssetAmountStat) TableName() string {
	return "asset_amount_stat"
}

type AssetAmountStatView interface {
	GetByGuid(guid string) (*AssetAmountStat, error)
	GetByAssetAndDate(assetUUID, date string) (*AssetAmountStat, error)
}

type AssetAmountStatDB interface {
	AssetAmountStatView

	StoreAssetAmountStat(a *AssetAmountStat) error
	StoreAssetAmountStats(list []*AssetAmountStat) error
	UpdateAssetAmountStat(guid string, updates map[string]interface{}) error
}

type assetAmountStatDB struct {
	gorm *gorm.DB
}

func NewAssetAmountStatDB(db *gorm.DB) AssetAmountStatDB {
	return &assetAmountStatDB{gorm: db}
}

func (db *assetAmountStatDB) StoreAssetAmountStat(a *AssetAmountStat) error {
	if err := db.gorm.Create(a).Error; err != nil {
		log.Error("StoreAssetAmountStat error", "err", err)
		return err
	}
	return nil
}

func (db *assetAmountStatDB) StoreAssetAmountStats(list []*AssetAmountStat) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreAssetAmountStats error", "err", err)
		return err
	}
	return nil
}

func (db *assetAmountStatDB) GetByGuid(guid string) (*AssetAmountStat, error) {
	var a AssetAmountStat
	if err := db.gorm.Where("guid = ?", guid).First(&a).Error; err != nil {
		log.Error("GetByGuid AssetAmountStat error", "err", err)
		return nil, err
	}
	return &a, nil
}

func (db *assetAmountStatDB) GetByAssetAndDate(assetUUID, date string) (*AssetAmountStat, error) {
	var a AssetAmountStat
	if err := db.gorm.Where("asset_uuid = ? AND time_date = ?", assetUUID, date).First(&a).Error; err != nil {
		log.Error("GetByAssetAndDate AssetAmountStat error", "err", err)
		return nil, err
	}
	return &a, nil
}

func (db *assetAmountStatDB) UpdateAssetAmountStat(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&AssetAmountStat{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateAssetAmountStat error", "err", err)
		return err
	}
	return nil
}
