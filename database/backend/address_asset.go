// address_asset.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AddressAsset struct {
	Guid        string          `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	TokenID     string          `gorm:"column:token_id;type:varchar(255);default:''" json:"token_id"`
	WalletUUID  string          `gorm:"column:wallet_uuid;type:varchar(255);default:''" json:"wallet_uuid"`
	AddressUUID string          `gorm:"column:address_uuid;type:varchar(255);default:''" json:"address_uuid"`
	AssetUsdt   decimal.Decimal `gorm:"column:asset_usdt;type:numeric(20,8);not null" json:"asset_usdt"`
	AssetUsd    decimal.Decimal `gorm:"column:asset_usd;type:numeric(20,8);not null" json:"asset_usd"`
	Balance     decimal.Decimal `gorm:"column:balance;type:numeric(78,0);not null" json:"balance"` // 使用 string 存储大数字（支持 uint256）
	CreateTime  time.Time       `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime  time.Time       `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (AddressAsset) TableName() string {
	return "address_asset"
}

type AddressAssetView interface {
	GetByGuid(guid string) (*AddressAsset, error)
	GetByAddressUUID(addressUUID string) ([]*AddressAsset, error)
	GetByAddressAndToken(addressUUID, tokenID string) (*AddressAsset, error)
}

type AddressAssetDB interface {
	AddressAssetView

	UpsertAddressAsset(a *AddressAsset) error
	UpsertAddressAssets(list []*AddressAsset) error
}

type addressAssetDB struct {
	db *gorm.DB
}

func NewAddressAssetDB(db *gorm.DB) AddressAssetDB {
	return &addressAssetDB{db: db}
}
func (r *addressAssetDB) UpsertAddressAsset(a *AddressAsset) error {
	if a.AddressUUID == "" || a.TokenID == "" {
		return fmt.Errorf("address_uuid and token_id cannot be empty")
	}

	a.UpdateTime = time.Now()

	err := r.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "address_uuid"}, {Name: "token_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"wallet_uuid",
				"asset_usdt",
				"asset_usd",
				"balance",
				"updated_at",
			}),
		}).
		Create(a).Error

	if err != nil {
		log.Error("UpsertAddressAsset error", "err", err)
		return err
	}
	return nil
}

func (r *addressAssetDB) UpsertAddressAssets(list []*AddressAsset) error {
	if len(list) == 0 {
		return nil
	}

	now := time.Now()
	for _, a := range list {
		if a.AddressUUID == "" || a.TokenID == "" {
			return fmt.Errorf("address_uuid and token_id cannot be empty")
		}
		a.UpdateTime = now
	}

	const batchSize = 100

	err := r.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "address_uuid"}, {Name: "token_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"wallet_uuid",
				"asset_usdt",
				"asset_usd",
				"balance",
				"updated_at",
			}),
		}).
		CreateInBatches(list, batchSize).Error

	if err != nil {
		log.Error("UpsertAddressAssets error", "err", err)
		return err
	}
	return nil
}

func (r *addressAssetDB) GetByGuid(guid string) (*AddressAsset, error) {
	var a AddressAsset
	if err := r.db.Where("guid = ?", guid).First(&a).Error; err != nil {
		log.Error("GetByGuid AddressAsset error", "err", err)
		return nil, err
	}
	return &a, nil
}

func (r *addressAssetDB) GetByAddressUUID(addressUUID string) ([]*AddressAsset, error) {
	var list []*AddressAsset
	if err := r.db.Where("address_uuid = ?", addressUUID).Find(&list).Error; err != nil {
		log.Error("GetByAddressUUID AddressAsset error", "err", err)
		return nil, err
	}
	return list, nil
}

func (r *addressAssetDB) GetByAddressAndToken(addressUUID, tokenID string) (*AddressAsset, error) {
	var a AddressAsset
	if err := r.db.Where("address_uuid = ? AND token_id = ?", addressUUID, tokenID).First(&a).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error("GetByAddressAndToken AddressAsset error", "err", err,
			"address_uuid", addressUUID, "token_id", tokenID)
		return nil, err
	}
	return &a, nil
}
