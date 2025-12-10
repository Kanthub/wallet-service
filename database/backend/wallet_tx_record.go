// wallet_tx_record.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

// TxStatus
const (
	TxStatusCreated = 0 // CREATED: 后端收到 signedTx 请求并写入记录，但尚未广播获得 txHash
	TxStatusPending = 1 // PENDING: 广播成功并拿到 txHash
	TxStatusFailed  = 2 // FAILED: 广播失败或链上执行失败或超时判失败
	TxStatusSuccess = 3 // SUCCESS: 链上确认成功
)

// TxStatus 状态名称映射
var TxStatusNames = map[int]string{
	TxStatusCreated: "CREATED",
	TxStatusPending: "PENDING",
	TxStatusFailed:  "FAILED",
	TxStatusSuccess: "SUCCESS",
}

// 失败原因代码常量
const (
	FailReasonBroadcastFailed = "BROADCAST_FAILED"  // 广播失败
	FailReasonChainFailed     = "CHAIN_FAILED"      // 链上执行失败
	FailReasonNotFoundTimeout = "NOT_FOUND_TIMEOUT" // 查不到且超时
	FailReasonUnknown         = "UNKNOWN"           // 未知错误
)

// 完整动作链路过滤：OperationID、StepIndex、TxType
type WalletTxRecord struct {
	Guid           string     `gorm:"primaryKey;column:guid;type:uuid;default:gen_random_uuid()" json:"guid"`
	OperationID    string     `gorm:"column:operation_id;type:varchar(255);default:'';index:idx_operation_step" json:"operation_id"` // 关联到完整操作（如 SwapID）
	StepIndex      int        `gorm:"column:step_index;type:integer;default:0;index:idx_operation_step" json:"step_index"`           // 步骤索引（0, 1, 2...）
	WalletUUID     string     `gorm:"column:wallet_uuid;type:varchar(255);not null;index" json:"wallet_uuid"`
	AddressUUID    string     `gorm:"column:address_uuid;type:varchar(255);default:'';index" json:"address_uuid"`
	TxTime         string     `gorm:"column:tx_time;type:varchar(500);not null" json:"tx_time"`
	ChainID        string     `gorm:"column:chain_id;type:varchar(255);default:'';index" json:"chain_id"`
	TokenID        string     `gorm:"column:token_id;type:varchar(255);default:''" json:"token_id"`
	FromAddress    string     `gorm:"column:from_address;type:varchar(70);not null;index" json:"from_address"`
	ToAddress      string     `gorm:"column:to_address;type:varchar(70);not null;index" json:"to_address"`
	Amount         string     `gorm:"column:amount;type:numeric(78,0);not null" json:"amount"` // 使用 string 存储大数字（支持 uint256）
	Memo           string     `gorm:"column:memo;type:varchar(500);not null" json:"memo"`
	Hash           string     `gorm:"column:hash;type:varchar(500);default:'';uniqueIndex" json:"hash"`
	BlockHeight    string     `gorm:"column:block_height;type:varchar(500);default:''" json:"block_height"`
	TxType         string     `gorm:"column:tx_type;type:varchar(50);default:'transfer';index" json:"tx_type"` // approve, swap, bridge, wrap, unwrap, transfer
	Status         int        `gorm:"column:status;type:integer;default:0;index:idx_status_last_checked" json:"status"`
	FailReasonCode string     `gorm:"column:fail_reason_code;type:varchar(100);default:''" json:"fail_reason_code,omitempty"`
	FailReasonMsg  string     `gorm:"column:fail_reason_msg;type:varchar(500);default:''" json:"fail_reason_msg,omitempty"`
	LastCheckedAt  *time.Time `gorm:"column:last_checked_at;index:idx_status_last_checked" json:"last_checked_at,omitempty"`
	CreateTime     time.Time  `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime     time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (WalletTxRecord) TableName() string {
	return "wallet_tx_record"
}

type WalletTxRecordView interface {
	GetByGuid(guid string) (*WalletTxRecord, error)
	GetByHash(hash string) (*WalletTxRecord, error)
	GetByOperationID(operationID string) ([]*WalletTxRecord, error)
	GetTxList(page, pageSize int, filters map[string]interface{}) ([]*WalletTxRecord, int64, error)
	GetPendingTxsForCheck(lastCheckedBefore time.Time, limit int) ([]*WalletTxRecord, error)
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

func (db *walletTxRecordDB) GetByOperationID(operationID string) ([]*WalletTxRecord, error) {
	var list []*WalletTxRecord
	if err := db.gorm.Where("operation_id = ?", operationID).Order("step_index ASC").Find(&list).Error; err != nil {
		log.Error("GetByOperationID WalletTxRecord error", "err", err)
		return nil, err
	}
	return list, nil
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
		case "from_address", "to_address", "hash", "wallet_uuid", "chain_id", "token_id", "tx_type":
			query = query.Where(key+" = ?", value)
		case "tx_status":
			query = query.Where("status = ?", value)
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

// GetPendingTxsForCheck 获取需要检查状态的 pending 交易
// lastCheckedBefore: 上次检查时间早于此时间的记录
// limit: 最多返回的记录数
func (db *walletTxRecordDB) GetPendingTxsForCheck(lastCheckedBefore time.Time, limit int) ([]*WalletTxRecord, error) {
	var list []*WalletTxRecord
	query := db.gorm.Model(&WalletTxRecord{}).
		Where("status = ?", TxStatusPending).
		Where("hash != ?", ""). // 必须有 hash
		Where("(last_checked_at IS NULL OR last_checked_at < ?)", lastCheckedBefore).
		Order("last_checked_at ASC NULLS FIRST"). // 优先处理从未检查过的
		Limit(limit)

	if err := query.Find(&list).Error; err != nil {
		log.Error("GetPendingTxsForCheck error", "err", err)
		return nil, err
	}

	return list, nil
}
