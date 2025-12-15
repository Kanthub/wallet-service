package backend

// 创建交易
type CreateWalletTxRequest struct {
	OperationID string `json:"operation_id"`
	StepIndex   int    `json:"step_index"`
	WalletUUID  string `json:"wallet_uuid"`
	AddressUUID string `json:"address_uuid"`
	TxTime      string `json:"tx_time"`
	ChainID     string `json:"chain_id"`
	TokenID     string `json:"token_id"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	Amount      string `json:"amount"`
	Memo        string `json:"memo"`
	TxType      string `json:"tx_type"`
}

type CreateWalletTxResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Guid    string `json:"guid"`
}

// 更新状态
type UpdateWalletTxStatusRequest struct {
	Guid           string `json:"guid"`
	Status         int    `json:"status"`
	Hash           string `json:"hash,omitempty"`
	BlockHeight    string `json:"block_height,omitempty"`
	FailReasonCode string `json:"fail_reason_code,omitempty"`
	FailReasonMsg  string `json:"fail_reason_msg,omitempty"`
}
