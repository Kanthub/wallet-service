# 多步骤操作设计文档

## 📋 概述

本文档描述了 wallet-services 中多步骤操作（Multi-Step Operation）的设计方案。该设计支持将一个完整的业务操作（如 Swap）拆分为多个交易步骤（如 Approve、Bridge、Swap、Wrap、Unwrap），并独立追踪每个步骤的状态。

---

## 🎯 设计目标

1. **通用性** - 不局限于 Swap，未来可支持 Transfer、Stake 等多步骤操作
2. **独立性** - 每条交易记录有独立的 UUID，便于追踪和查询
3. **关联性** - 通过 `operation_id` 关联所有相关交易
4. **可扩展** - 新增操作类型只需添加常量定义
5. **完整性** - 保存所有步骤（approve/bridge/swap/wrap/unwrap）
6. **可查询** - 可按钱包、操作、步骤等多维度查询

---

## 🏗️ 核心概念

```
Operation (操作) - 一个完整的业务操作
  ├── OperationID (操作ID) - 唯一标识一个完整操作（如 SwapID）
  ├── OperationType (操作类型) - "SWAP", "TRANSFER", "STAKE" 等
  └── Steps (步骤列表) - 包含多个交易步骤
       ├── Step 0: APPROVE (TxRecord with unique GUID)
       ├── Step 1: BRIDGE (TxRecord with unique GUID)
       └── Step 2: SWAP (TxRecord with unique GUID)
```

---

## 📊 数据模型

### WalletTxRecord 表结构

```go
type WalletTxRecord struct {
    Guid           string     // UUID - 每条记录的唯一ID
    OperationID    string     // 关联到完整操作（如 SwapID）
    StepIndex      int        // 步骤索引（0, 1, 2...）
    WalletUUID     string     // 钱包UUID
    TxType         string     // 交易类型：approve, swap, bridge, wrap, unwrap, transfer
    Status         int        // 交易状态：0=CREATED, 1=PENDING, 2=FAILED, 3=SUCCESS
    Hash           string     // 交易哈希
    // ... 其他字段
}
```

**关键字段说明：**

- **Guid** - 使用 UUID 作为唯一主键（每条交易记录独立的 ID）
- **OperationID** - 关联到完整操作（原来的 SwapID），用于查询某个操作的所有步骤
- **StepIndex** - 记录在操作中的步骤顺序（0, 1, 2...）
- **TxType** - 记录具体的交易类型（approve/swap/bridge/wrap/unwrap/transfer）

**索引：**
- `idx_operation_step` - 复合索引 (operation_id, step_index)，用于高效查询操作的所有步骤

---

## 🔄 数据流

### 1. 创建操作

```go
// 生成 OperationID（SwapID）
operationID := uuid.New().String()

// 为每个 step 创建独立的 WalletTxRecord
for i, action := range actions {
    recordGuid := uuid.New().String() // 每条记录独立的 UUID
    
    record := &WalletTxRecord{
        Guid:        recordGuid,      // 独立的 UUID
        OperationID: operationID,     // 关联到操作
        StepIndex:   i,               // 步骤索引
        TxType:      string(action.ActionType), // approve, swap, bridge, wrap, unwrap
        Status:      TxStatusCreated,
        // ...
    }
    
    db.Create(record)
}
```

### 2. 查询操作的所有步骤

```go
// 查询某个操作的所有交易记录（按步骤顺序）
var records []WalletTxRecord
db.Where("operation_id = ?", operationID).
   Order("step_index ASC").
   Find(&records)
```

### 3. 查询钱包的所有交易

```go
// 查询某个钱包的所有交易记录（按时间倒序）
var records []WalletTxRecord
db.Where("wallet_uuid = ?", walletUUID).
   Order("tx_time DESC").
   Find(&records)
```

---

## 🎨 API 响应示例

```json
{
  "swap_id": "550e8400-e29b-41d4-a716-446655440000",
  "wallet_uuid": "wallet-123",
  "user_address": "0x1234...",
  "state": "PENDING",
  "steps": [
    {
      "step_index": 0,
      "action_type": "APPROVE",
      "state": "CONFIRMED",
      "tx_hash": "0xabc123...",
      "tx_record_guid": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
    },
    {
      "step_index": 1,
      "action_type": "BRIDGE",
      "state": "PENDING",
      "tx_hash": "0xdef456...",
      "tx_record_guid": "b2c3d4e5-f6a7-8901-bcde-f12345678901"
    },
    {
      "step_index": 2,
      "action_type": "SWAP",
      "state": "PENDING",
      "tx_hash": "",
      "tx_record_guid": null
    }
  ]
}
```

---

## 🔧 ActionType 定义

```go
type ActionType string

const (
    ActionTypeApprove ActionType = "APPROVE" // 授权操作
    ActionTypeSwap    ActionType = "SWAP"    // 交换操作
    ActionTypeBridge  ActionType = "BRIDGE"  // 跨链桥接
    ActionTypeWrap    ActionType = "WRAP"    // 包装原生代币（ETH -> WETH）
    ActionTypeUnwrap  ActionType = "UNWRAP"  // 解包装代币（WETH -> ETH）
)
```

---

## 📝 数据库迁移

迁移脚本：`database/migrations/008_add_operation_fields_to_wallet_tx_record.sql`

```sql
-- Add operation_id field
ALTER TABLE wallet_tx_record 
ADD COLUMN IF NOT EXISTS operation_id VARCHAR(255) DEFAULT '' NOT NULL;

-- Add step_index field
ALTER TABLE wallet_tx_record 
ADD COLUMN IF NOT EXISTS step_index INTEGER DEFAULT 0 NOT NULL;

-- Create composite index
CREATE INDEX IF NOT EXISTS idx_operation_step ON wallet_tx_record(operation_id, step_index);
```

---

## ✅ 优势

1. **完整的交易历史** - 用户可以看到所有步骤（approve + bridge + swap）
2. **独立追踪** - 每个步骤有独立的 UUID 和状态
3. **灵活查询** - 可按操作、钱包、步骤等多维度查询
4. **易于扩展** - 未来可支持更多操作类型（Transfer、Stake 等）
5. **故障排查** - 可以精确定位哪一步失败
6. **Worker 兼容** - 现有的 Worker 无需修改，自动支持所有交易类型

---

## 🚀 未来扩展

### 支持更多操作类型

```go
type OperationType string

const (
    OperationTypeSwap     OperationType = "SWAP"
    OperationTypeTransfer OperationType = "TRANSFER"
    OperationTypeStake    OperationType = "STAKE"
    // 未来可以扩展更多类型
)
```

### 操作级别的状态追踪

可以考虑添加独立的 `operations` 表来追踪操作级别的状态：

```go
type Operation struct {
    OperationID    string
    OperationType  string
    WalletUUID     string
    State          string // PENDING, SUCCESS, FAILED
    TotalSteps     int
    CompletedSteps int
}
```

---

## 📚 相关文件

- `database/backend/wallet_tx_record.go` - 数据模型定义
- `services/api/service/aggregator_service.go` - 业务逻辑实现
- `services/api/models/backend/aggregator.go` - API 模型定义
- `worker/aggregator_task/wallet_tx_record_worker.go` - 后台 Worker
- `database/migrations/008_add_operation_fields_to_wallet_tx_record.sql` - 数据库迁移脚本

