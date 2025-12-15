package routes

import (
	"encoding/json"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"
	"github.com/roothash-pay/wallet-services/rpc/account"
	"github.com/roothash-pay/wallet-services/services/api/service"
)

type TransactionInfo struct {
	TxHash string `json:"tx_hash"`
	RawTx  string `json:"raw_tx"`
}

func (rs *Routes) ChainInfoApi() {
	r := rs.router

	r.Route("/api/v1/chain_info", func(r chi.Router) {
		r.Get("/balance", rs.getChainBalance)
		r.Get("/balance_sync", rs.getChainBalanceSync)
		r.Get("/sign_info", rs.getSignInfo)
		r.Get("/send_tx", rs.sendRawTransaction)
		r.Get("/transactions", rs.getTransactions)
		r.Post("/submit_tx", rs.submitTx)
		r.Get("/txn_status", rs.getTxnStatus)
	})
}

func (rs *Routes) getChainBalanceSync(w http.ResponseWriter, r *http.Request) {

	chainStr := r.URL.Query().Get("chain")
	if chainStr == "" {
		chainStr = "polygon" // 默认链
	}

	tokenAddress := r.URL.Query().Get("token_address")

	chainType := service.ChainType(chainStr)

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address required", http.StatusBadRequest)
		return
	}

	nativeBalance, err := rs.svc.DappLinkService.GetBalanceByAddress(chainType, address)
	if err != nil {
		log.Error("get native balance error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	erc20Balance, err := rs.svc.DappLinkService.GetErc20BalanceByAddress(chainType,
		tokenAddress, address,
	)
	if err != nil {
		log.Error("get usdt balance error", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{
		"native_balance": nativeBalance,
		"erc20_balance":  erc20Balance,
	}, http.StatusOK)
}

func (rs *Routes) getChainBalance(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address required", http.StatusBadRequest)
		return
	}

	tokenAddress := r.URL.Query().Get("token_address")

	chain := r.URL.Query().Get("chain")
	if chain == "" {
		chain = "roothash" // 默认链
	}

	ctx := r.Context()

	get := func(contract string) (*account.AccountResponse, error) {
		return rs.svc.RpcService.GetAccount(ctx, &account.AccountRequest{
			Chain:           chain,
			Network:         "mainnet",
			Address:         address,
			ContractAddress: contract,
		})
	}

	nativeBalance, err := get("0x00")
	if err != nil || nativeBalance.Code == account.ReturnCode_ERROR {
		http.Error(w, "get native balance error", http.StatusInternalServerError)
		return
	}

	erc20Balance, err := get(tokenAddress)
	if err != nil || erc20Balance.Code == account.ReturnCode_ERROR {
		http.Error(w, "get erc20 balance error", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{
		"native_balance": nativeBalance.Balance,
		"erc20_balance":  erc20Balance.Balance,
	}, http.StatusOK)
}

func (rs *Routes) getSignInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address required", http.StatusBadRequest)
		return
	}

	chain := r.URL.Query().Get("chain")
	if chain == "" {
		chain = "roothash"
	}

	// 1️⃣ 获取账户信息（nonce）
	resp, err := rs.svc.RpcService.GetAccount(
		ctx,
		&account.AccountRequest{
			Chain:   chain,
			Network: "mainnet",
			Address: address,
		},
	)
	if err != nil || resp.Code == account.ReturnCode_ERROR {
		http.Error(w, "get account error", http.StatusInternalServerError)
		return
	}

	// 2️⃣ feeHistory
	const blockCount = 4
	percentiles := []int{70, 95}

	var feeHist struct {
		OldestBlock   string     `json:"oldestBlock"`
		BaseFeePerGas []string   `json:"baseFeePerGas"`
		GasUsedRatio  []float64  `json:"gasUsedRatio"`
		Reward        [][]string `json:"reward"`
	}

	client, ok := rs.svc.Client[service.ChainType(chain)]
	if !ok {
		http.Error(w, "unsupported chain", http.StatusBadRequest)
		return
	}

	if err := client.CallContext(
		ctx,
		&feeHist,
		"eth_feeHistory",
		blockCount,
		"latest",
		percentiles,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(feeHist.BaseFeePerGas) == 0 {
		http.Error(w, "invalid baseFee data", http.StatusInternalServerError)
		return
	}

	latestBaseFeeHex := feeHist.BaseFeePerGas[len(feeHist.BaseFeePerGas)-1]
	latestBaseFee, ok := new(big.Int).SetString(strings.TrimPrefix(latestBaseFeeHex, "0x"), 16)
	if !ok {
		http.Error(w, "baseFee parse error", http.StatusInternalServerError)
		return
	}

	if len(feeHist.Reward) == 0 {
		http.Error(w, "invalid reward data", http.StatusInternalServerError)
		return
	}

	lastReward := feeHist.Reward[len(feeHist.Reward)-1]
	if len(lastReward) < 2 {
		http.Error(w, "reward percentile missing", http.StatusInternalServerError)
		return
	}

	tip70, ok := new(big.Int).SetString(strings.TrimPrefix(lastReward[0], "0x"), 16)
	if !ok {
		http.Error(w, "tip70 parse error", http.StatusInternalServerError)
		return
	}

	tip95, ok := new(big.Int).SetString(strings.TrimPrefix(lastReward[1], "0x"), 16)
	if !ok {
		http.Error(w, "tip95 parse error", http.StatusInternalServerError)
		return
	}

	// 3️⃣ gas 组合
	slow := new(big.Int).Add(latestBaseFee, tip70)
	normal := new(big.Int).Add(latestBaseFee, tip95)

	retValue := map[string]interface{}{
		"nonce":                    resp.Sequence,
		"native_token_gas_limit":   "21000",
		"erc20_token_gas_limit":    "150000",
		"max_fee_per_gas":          normal.String(),
		"max_priority_fee_per_gas": tip95.String(),
		"gas_price":                slow.String(), // legacy 兼容
	}

	jsonResponse(w, retValue, http.StatusOK)
}

func (rs *Routes) sendRawTransaction(w http.ResponseWriter, r *http.Request) {
	rawTx := r.URL.Query().Get("rawTx")
	if rawTx == "" {
		http.Error(w, "rawTx required", http.StatusBadRequest)
		return
	}

	resp, err := rs.svc.RpcService.SendTx(
		r.Context(),
		&account.SendTxRequest{
			Chain:   "Polygon",
			Network: "mainnet",
			RawTx:   rawTx,
		},
	)
	if err != nil || resp == nil || resp.Code == account.ReturnCode_ERROR {
		log.Error("send raw tx error", "err", err)
		http.Error(w, "send tx failed", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{
		"tx_hash": resp.TxHash,
	}, http.StatusOK)
}

func (rs *Routes) getTransactions(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "address required", http.StatusBadRequest)
		return
	}

	chain := r.URL.Query().Get("chain")
	if chain == "" {
		chain = "roothash"
	}

	resp, err := rs.svc.RpcService.GetTxByAddress(
		r.Context(),
		&account.TxAddressRequest{
			Chain:   chain,
			Network: "mainnet",
			Address: address,
		},
	)
	if err != nil || resp.Code == account.ReturnCode_ERROR {
		http.Error(w, "get transactions error", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, resp, http.StatusOK)
}

func (rs *Routes) submitTx(w http.ResponseWriter, r *http.Request) {
	var req TransactionInfo
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if rs.svc.WalletService.IsExistRawTx(req.RawTx) {
		http.Error(w, "raw tx already exist", http.StatusBadRequest)
		return
	}

	if err := rs.svc.WalletService.StoreRawTx(req.RawTx, req.TxHash); err != nil {
		log.Error("store raw tx error", "err", err)
		http.Error(w, "store tx failed", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, "ok", http.StatusOK)
}

func (rs *Routes) getTxnStatus(w http.ResponseWriter, r *http.Request) {
	txHash := r.URL.Query().Get("hash")
	if txHash == "" {
		http.Error(w, "hash required", http.StatusBadRequest)
		return
	}

	tx, err := rs.svc.WalletService.QueryTxInfoByHash(txHash)
	if err != nil {
		log.Error("query tx status error", "err", err)
		http.Error(w, "query tx failed", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, tx, http.StatusOK)
}
