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
	ActionTypeApprove ActionType = "APPROVE"
	ActionTypeSwap    ActionType = "SWAP"
)

// SwapState represents the state of a swap
type SwapState string

const (
	SwapStatePending   SwapState = "PENDING"
	SwapStateSubmitted SwapState = "SUBMITTED"
	SwapStateConfirmed SwapState = "CONFIRMED"
	SwapStateFailed    SwapState = "FAILED"
)

// StepState represents the state of a transaction step
type StepState string

const (
	StepStatePending   StepState = "PENDING"
	StepStateSubmitted StepState = "SUBMITTED"
	StepStateConfirmed StepState = "CONFIRMED"
	StepStateFailed    StepState = "FAILED"
)

// QuoteRequest represents a request for swap quotes
type QuoteRequest struct {
	FromChainID string `json:"from_chain_id" validate:"required"`
	ToChainID   string `json:"to_chain_id" validate:"required"`
	FromToken   string `json:"from_token" validate:"required"`
	ToToken     string `json:"to_token" validate:"required"`
	Amount      string `json:"amount" validate:"required"`
	SlippageBps int    `json:"slippage_bps" validate:"required,min=0,max=10000"`
	UserAddress string `json:"user_address,omitempty"`
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

// QuoteResponse represents the response containing quotes
type QuoteResponse struct {
	QuoteID      string    `json:"quote_id"`
	ExpiresAt    time.Time `json:"expires_at"`
	BestQuote    *Quote    `json:"best_quote"`
	Alternatives []*Quote  `json:"alternatives,omitempty"`
}

// PrepareSwapRequest represents a request to prepare a swap
type PrepareSwapRequest struct {
	QuoteID     string `json:"quote_id" validate:"required"`
	UserAddress string `json:"user_address" validate:"required"`
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

// Step represents a single transaction step in a swap
type Step struct {
	StepIndex      int        `json:"step_index"`
	ActionType     ActionType `json:"action_type"`
	TxHash         string     `json:"tx_hash,omitempty"`
	State          StepState  `json:"state"`
	SubmittedAt    *time.Time `json:"submitted_at,omitempty"`
	ConfirmedAt    *time.Time `json:"confirmed_at,omitempty"`
	FailReasonCode string     `json:"fail_reason_code,omitempty"`
	FailMessage    string     `json:"fail_message,omitempty"`
	IdempotencyKey string     `json:"idempotency_key,omitempty"`
}

// Swap represents a complete swap operation
type Swap struct {
	SwapID         string    `json:"swap_id"`
	QuoteID        string    `json:"quote_id"`
	UserAddress    string    `json:"user_address"`
	State          SwapState `json:"state"`
	Steps          []*Step   `json:"steps"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	FailReasonCode string    `json:"fail_reason_code,omitempty"`
	FailMessage    string    `json:"fail_message,omitempty"`
}

// SwapStatusResponse represents the response for swap status query
type SwapStatusResponse struct {
	SwapID         string    `json:"swap_id"`
	State          SwapState `json:"state"`
	Steps          []*Step   `json:"steps"`
	FailReasonCode string    `json:"fail_reason_code,omitempty"`
	FailMessage    string    `json:"fail_message,omitempty"`
}
