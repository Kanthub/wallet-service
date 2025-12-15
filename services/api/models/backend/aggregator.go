package backend

import (
	"time"
)

// ChainType represents the blockchain type
type ChainType string

const (
	ChainTypeEVM    ChainType = "EVM"
	ChainTypeSolana ChainType = "SOLANA"
)

// ActionType represents the type of action in a transaction plan
type ActionType string

const (
	ActionTypeApprove ActionType = "APPROVE" // 授权操作
	ActionTypeSwap    ActionType = "SWAP"    // 交换操作
	ActionTypeBridge  ActionType = "BRIDGE"  // 跨链桥接
	ActionTypeWrap    ActionType = "WRAP"    // 包装原生代币
	ActionTypeUnwrap  ActionType = "UNWRAP"  // 解包装代币
)

// TxStatus represents the status of a transaction (unified with database)
// 使用与 database/backend/wallet_tx_record.go 相同的状态定义
const (
	TxStatusCreated = 0 // CREATED: 后端收到 signedTx 请求并写入记录，但尚未广播
	TxStatusPending = 1 // PENDING: 广播成功并拿到 txHash
	TxStatusFailed  = 2 // FAILED: 广播失败或链上执行失败或超时
	TxStatusSuccess = 3 // SUCCESS: 链上确认成功
)

// TxStatusNames provides human-readable names for status codes
var TxStatusNames = map[int]string{
	TxStatusCreated: "CREATED",
	TxStatusPending: "PENDING",
	TxStatusFailed:  "FAILED",
	TxStatusSuccess: "SUCCESS",
}

type RoutesRequest struct {
	FromChainId      string               `json:"fromChainId"`
	FromAmount       string               `json:"fromAmount"`
	FromTokenAddress string               `json:"fromTokenAddress"`
	ToChainId        string               `json:"toChainId"`
	ToTokenAddress   string               `json:"toTokenAddress"`
	FromAddress      string               `json:"fromAddress,omitempty"` // 可选用户地址
	Slippage         float64              `json:"slippage,omitempty"`    // 滑点
	Options          RoutesRequestOptions `json:"options"`
}

// RoutesRequestOptions routes请求的options参数
type RoutesRequestOptions struct {
	Bridges              RoutesAllowList `json:"bridges"`
	Exchanges            RoutesAllowList `json:"exchanges"`
	AllowSwitchChain     bool            `json:"allowSwitchChain"`
	AllowDestinationCall bool            `json:"allowDestinationCall"`
}

// RoutesAllowList 允许的桥/交易所列表
type RoutesAllowList struct {
	Allow []string `json:"allow"`
}

// QuoteRequest represents a request for swap quotes
type QuoteRequest struct {
	FromChainID string  `json:"from_chain_id" validate:"required"`
	ToChainID   string  `json:"to_chain_id" validate:"required"`
	FromToken   string  `json:"from_token" validate:"required"`
	ToToken     string  `json:"to_token" validate:"required"`
	Amount      string  `json:"amount" validate:"required"`
	SlippageBps float64 `json:"slippage_bps" validate:"required,min=0,max=10000"`
	UserAddress string  `json:"user_address,omitempty"`
	WalletUUID  string  `json:"wallet_uuid,omitempty"` // Optional: wallet UUID for tracking
}

// Quote represents a swap quote from a provider
type Quote struct {
	Provider    string    `json:"provider"`
	ChainType   ChainType `json:"chain_type"`
	ChainID     string    `json:"chain_id"`
	FromToken   string    `json:"from_token"`
	ToToken     string    `json:"to_token"`
	FromAmount  string    `json:"from_amount"`
	ToAmount    string    `json:"to_amount"`
	GasEstimate string    `json:"gas_estimate"`
	Fees        string    `json:"fees,omitempty"`
	Spender     string    `json:"spender,omitempty"` // EVM only: approval spender
	Router      string    `json:"router,omitempty"`  // EVM only: swap router
	Raw         string    `json:"raw,omitempty"`     // Raw provider response
}

type QuoteStore struct {
	QuoteID         string    `json:"quote_id"`
	UserAddress     string    `json:"user_address"`
	WalletUUID      string    `json:"wallet_uuid"`
	ExpiresAt       time.Time `json:"expires_at"`
	BestQuotesIndex int       `json:"best_quotes_index"`
	BestQuotes      []*Quote  `json:"best_quotes"`
	Raws            []string  `json:"raws,omitempty"`
}

// QuoteResponse represents the response containing quotes
type QuoteResponse struct {
	QuoteID     string    `json:"quote_id"`
	UserAddress string    `json:"user_address"`
	WalletUUID  string    `json:"wallet_uuid"`
	ExpiresAt   time.Time `json:"expires_at"`
	BestQuotes  []*Quote  `json:"best_quotes"`
}

// PrepareSwapRequest represents a request to prepare a swap
type PrepareSwapRequest struct {
	QuoteID         string `json:"quote_id" validate:"required"`
	BestQuotesIndex int    `json:"best_quotes_index" validate:"required"`
}

// SigningPayload represents the data to be signed
type SigningPayload struct {
	// EVM fields
	To      string `json:"to,omitempty"`
	Data    string `json:"data,omitempty"`
	Value   string `json:"value,omitempty"`
	Gas     string `json:"gas,omitempty"`
	ChainID string `json:"chain_id,omitempty"`

	// Solana fields
	SerializedTx string `json:"serialized_tx,omitempty"` // TODO: Solana transaction payload
}

// Action represents a single action in a transaction plan
type Action struct {
	ActionType     ActionType      `json:"action_type"`
	ChainID        string          `json:"chain_id"`
	SigningPayload *SigningPayload `json:"signing_payload"`
	Description    string          `json:"description,omitempty"`
}

// TxPlan represents a plan of actions to execute a swap
type TxPlan struct {
	SwapID  string    `json:"swap_id"`
	Actions []*Action `json:"actions"`
}

// PrepareSwapResponse represents the response from prepare swap
type PrepareSwapResponse struct {
	SwapID  string    `json:"swap_id"`
	Actions []*Action `json:"actions"`
}

// BuildSwapResponse represents the response from provider build swap
type BuildSwapResponse struct {
	Actions []*Action `json:"actions"`
}

// SubmitSignedTxRequest represents a request to submit a signed transaction
type SubmitSignedTxRequest struct {
	SwapID         string `json:"swap_id" validate:"required"`
	StepIndex      int    `json:"step_index" validate:"required,min=0"`
	SignedTx       string `json:"signed_tx" validate:"required"`
	IdempotencyKey string `json:"idempotency_key" validate:"required"`
}

// SubmitSignedTxResponse represents the response from submitting a signed tx
type SubmitSignedTxResponse struct {
	TxHash string `json:"tx_hash"`
}

type SubmitTxHashRequest struct {
	SwapID         string `json:"swap_id" validate:"required"`
	StepIndex      int    `json:"step_index" validate:"min=0"`
	TxHash         string `json:"tx_hash" validate:"required"`
	IdempotencyKey string `json:"idempotency_key" validate:"required"`
}

type SubmitTxHashResponse struct {
	TxHash string `json:"tx_hash"`
}

// Step represents a single transaction step in a swap
type Step struct {
	StepIndex        int        `json:"step_index"`
	ActionType       ActionType `json:"action_type"`
	TxHash           string     `json:"tx_hash,omitempty"`
	Status           int        `json:"status"` // 0=CREATED, 1=PENDING, 2=FAILED, 3=SUCCESS
	SubmittedAt      *time.Time `json:"submitted_at,omitempty"`
	ConfirmedAt      *time.Time `json:"confirmed_at,omitempty"`
	FailReasonCode   string     `json:"fail_reason_code,omitempty"`
	FailMessage      string     `json:"fail_message,omitempty"`
	IdempotencyKey   string     `json:"idempotency_key,omitempty"`
	ExpectedChainID  string     `json:"expected_chain_id,omitempty"` // 用 string 存更通用
	ExpectedTo       string     `json:"expected_to,omitempty"`
	ExpectedValueWei string     `json:"expected_value,omitempty"`     // wei，十进制或 hex 统一一种
	ExpectedDataHash string     `json:"expected_data_hash,omitempty"` // 0x...
}

// Swap represents a complete swap operation
type Swap struct {
	SwapID         string    `json:"swap_id"`
	QuoteID        string    `json:"quote_id"` // 待修改
	UserAddress    string    `json:"user_address"`
	WalletUUID     string    `json:"wallet_uuid,omitempty"` // Wallet UUID for tracking
	Status         int       `json:"status"`                // 整体状态（根据所有 steps 计算）: 0=CREATED, 1=PENDING, 2=FAILED, 3=SUCCESS
	Steps          []*Step   `json:"steps"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	FailReasonCode string    `json:"fail_reason_code,omitempty"`
	FailMessage    string    `json:"fail_message,omitempty"`
}

// SwapStatusResponse represents the response for swap status query
type SwapStatusResponse struct {
	SwapID         string  `json:"swap_id"`
	Status         int     `json:"status"` // 0=CREATED, 1=PENDING, 2=FAILED, 3=SUCCESS
	Steps          []*Step `json:"steps"`
	FailReasonCode string  `json:"fail_reason_code,omitempty"`
	FailMessage    string  `json:"fail_message,omitempty"`
}
