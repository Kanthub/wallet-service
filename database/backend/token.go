package backend

import (
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type Token struct {
	Guid                 string    `gorm:"primaryKey;column:guid;type:text" json:"guid"`
	TokenName            string    `gorm:"column:token_name;type:varchar(70);default:''" json:"token_name"`
	TokenMark            string    `gorm:"column:token_mark;type:varchar(70);default:''" json:"token_mark"`
	TokenLogo            string    `gorm:"column:token_logo;type:varchar(100);default:''" json:"token_logo"`
	TokenActiveLogo      string    `gorm:"column:token_active_logo;type:varchar(100);default:''" json:"token_active_logo"`
	TokenDecimal         string    `gorm:"column:token_decimal;type:varchar(10);default:'18'" json:"token_decimal"`
	TokenSymbol          string    `gorm:"column:token_symbol;type:varchar(70);default:''" json:"token_symbol"`
	TokenContractAddress string    `gorm:"column:token_contract_address;type:varchar(70);not null" json:"token_contract_address"`
	ChainID              string    `gorm:"column:token_chain_id;type:varchar(255);default:''" json:"chain_id"`
	IsHot                string    `gorm:"column:is_hot;type:varchar(32);not null;default:'hot'" json:"is_hot"`
	CreateTime           time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime           time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (Token) TableName() string {
	return "token"
}

type TokenView interface {
	GetByGuid(guid string) (*Token, error)
	GetByContractAddress(addr string) (*Token, error)
	GetByContractAndChain(addr, chainID string) (*Token, error)
	GetTokenList(page, pageSize int, filters map[string]interface{}) ([]*Token, int64, error)
}

type TokenDB interface {
	TokenView

	StoreToken(t *Token) error
	StoreTokens(list []*Token) error
	UpdateToken(guid string, updates map[string]interface{}) error
}

type tokenDB struct {
	gorm *gorm.DB
}

func NewTokenDB(db *gorm.DB) TokenDB {
	return &tokenDB{gorm: db}
}

func (db *tokenDB) StoreToken(t *Token) error {
	if err := db.gorm.Create(t).Error; err != nil {
		log.Error("StoreToken error", "err", err)
		return err
	}
	return nil
}

func (db *tokenDB) StoreTokens(list []*Token) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreTokens error", "err", err)
		return err
	}
	return nil
}

func (db *tokenDB) GetByGuid(guid string) (*Token, error) {
	var t Token
	if err := db.gorm.Where("guid = ?", guid).First(&t).Error; err != nil {
		log.Error("GetByGuid token error", "err", err)
		return nil, err
	}
	return &t, nil
}

func (db *tokenDB) GetByContractAddress(addr string) (*Token, error) {
	var t Token
	if err := db.gorm.Where("token_contract_address = ?", addr).First(&t).Error; err != nil {
		log.Error("GetByContractAddress error", "err", err)
		return nil, err
	}
	return &t, nil
}

func (db *tokenDB) GetByContractAndChain(addr, chainID string) (*Token, error) {
	if addr == "" {
		return nil, fmt.Errorf("contract address is required")
	}

	normalized := strings.ToLower(addr)

	query := db.gorm.Where("LOWER(token_contract_address) = ?", normalized)
	if chainID != "" {
		query = query.Where("token_chain_id = ?", chainID)
	}

	var t Token
	if err := query.First(&t).Error; err != nil {
		log.Error("GetByContractAndChain token error", "addr", normalized, "chain_id", chainID, "err", err)
		return nil, err
	}
	return &t, nil
}

func (db *tokenDB) GetTokenList(page, pageSize int, filters map[string]interface{}) ([]*Token, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var list []*Token
	query := db.gorm.Model(&Token{})

	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}
		switch key {
		case "token_name", "token_symbol":
			query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error("GetTokenList count error", "err", err)
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
		log.Error("GetTokenList list error", "err", err)
		return nil, 0, err
	}

	return list, total, nil
}

func (db *tokenDB) UpdateToken(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}
	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&Token{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateToken error", "err", err)
		return err
	}
	return nil
}
