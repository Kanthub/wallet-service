# 数据库数值类型指南

## 📋 问题背景

在区块链应用中，使用 `INTEGER` 类型存储金额是**严重的设计缺陷**。

### INTEGER 的限制

PostgreSQL 的 `INTEGER` 类型范围：
- **最小值**: `-2,147,483,648`
- **最大值**: `2,147,483,647` (约 2.1 × 10^9)

### 区块链金额的特点

以太坊为例：
- **ETH** 使用 18 位小数，1 ETH = 10^18 wei
- **USDT/USDC** 使用 6 位小数，1 USDT = 10^6 最小单位
- 存储时保存**最小单位的整数值**

**示例：**
```
1 ETH = 1,000,000,000,000,000,000 wei (10^18)
这个数字远超 INTEGER 的最大值 2.1 × 10^9
```

---

## ✅ 解决方案

### 使用 NUMERIC 类型

```sql
-- 对于区块链金额（支持 uint256）
NUMERIC(78, 0)  -- 可存储 0 到 10^78-1 的整数

-- 对于美元金额（市值、流动性等）
NUMERIC(20, 2)  -- 可存储 0 到 10^18 的数字，精度到分

-- 对于价格（USDT/USD 价格）
NUMERIC(20, 8)  -- 可存储 0 到 10^12 的数字，精度到 0.00000001
```

---

## 🔧 修改的表

### 1. wallet_tx_record（交易记录）
```sql
amount NUMERIC(78,0) NOT NULL CHECK (amount >= 0)
```

### 2. wallet_asset（钱包资产）
```sql
balance NUMERIC(78,0) NOT NULL CHECK (balance >= 0)
```

### 3. asset_amount_stat（资产统计）
```sql
amount NUMERIC(78,0) NOT NULL CHECK (amount >= 0)
```

### 4. address_asset（地址资产）
```sql
balance NUMERIC(78,0) NOT NULL CHECK (balance >= 0)
```

### 5. market_price（市场价格）
```sql
market_cap   NUMERIC(20, 2) CHECK (market_cap >= 0)      -- 市值（美元）
liquidity    NUMERIC(20, 2) CHECK (liquidity >= 0)       -- 流动性（美元）
24h_volume   NUMERIC(20, 2) CHECK (24h_volume >= 0)      -- 24小时成交量（美元）
```

---

## 📊 NUMERIC 类型选择指南

| 用途 | 类型 | 范围 | 精度 | 说明 |
|------|------|------|------|------|
| 区块链金额 | NUMERIC(78,0) | 0 ~ 10^78-1 | 整数 | 支持 uint256 |
| 美元金额 | NUMERIC(20,2) | 0 ~ 10^18 | 0.01 | 市值、流动性 |
| 价格 | NUMERIC(20,8) | 0 ~ 10^12 | 0.00000001 | USDT/USD 价格 |
| 百分比 | NUMERIC(5,2) | 0 ~ 999.99 | 0.01 | 涨跌幅 |

---

## ⚠️ 注意事项

1. **CHECK 约束** - 改为 `>= 0` 而不是 `> 0`，允许零值
2. **性能** - NUMERIC 比 INTEGER 稍慢，但精度和范围更重要
3. **应用层** - Go 中使用 `decimal.Decimal` 或 `big.Int` 处理
4. **JSON 序列化** - 确保正确序列化大数字（避免精度丢失）

---

## 🔍 验证

### 数据库 Schema 修改
所有修改已应用到 `migrations/20251117001.sql`：
- ✅ wallet_tx_record.amount → NUMERIC(78,0)
- ✅ wallet_asset.balance → NUMERIC(78,0)
- ✅ asset_amount_stat.amount → NUMERIC(78,0)
- ✅ address_asset.balance → NUMERIC(78,0)
- ✅ market_price.market_cap → NUMERIC(20,2)
- ✅ market_price.liquidity → NUMERIC(20,2)
- ✅ market_price.24h_volume → NUMERIC(20,2)

### Go 结构体修改
所有结构体已更新为使用 `string` 类型：
- ✅ `database/backend/wallet_tx_record.go` - Amount: int64 → string
- ✅ `database/backend/wallet_asset.go` - Balance: int64 → string
- ✅ `database/backend/asset_amount_stat.go` - Amount: int64 → string
- ✅ `database/backend/address_asset.go` - Balance: int64 → string
- ✅ `database/backend/market_price.go` - MarketCap/Liquidity/Volume24h: int64 → string

### 业务逻辑修改
- ✅ `services/api/service/aggregator_service.go` - 直接使用字符串存储金额，无需转换

---

## 💡 Go 代码处理建议

### 使用 big.Int 处理大数字

```go
import "math/big"

// 从字符串解析
amount := new(big.Int)
amount.SetString("1000000000000000000", 10) // 1 ETH in wei

// 转换为字符串存储到数据库
amountStr := amount.String()

// 数学运算
balance := new(big.Int).SetString("5000000000000000000", 10)
result := new(big.Int).Add(amount, balance)
```

### JSON 序列化注意事项

```go
type Response struct {
    Amount string `json:"amount"` // 使用 string 避免精度丢失
}

// ❌ 错误：使用 int64 会导致精度丢失
// Amount int64 `json:"amount"`

// ✅ 正确：使用 string 保持精度
// Amount string `json:"amount"`
```

