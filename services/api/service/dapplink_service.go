package service

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/chain-explorer-api/common/chain"
	"github.com/dapplink-labs/chain-explorer-api/explorer/etherscan"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/log"

	"github.com/roothash-pay/wallet-services/config"
)

type ChainType string

const (
	ChainETH      ChainType = "eth"
	ChainArbitrum ChainType = "arbitrum"
	ChainPolygon  ChainType = "polygon"
	ChainBase     ChainType = "base"
	ChainOp       ChainType = "op"
	ChainRootHash ChainType = "roothash"
)

type DappLinkService interface {
	GetErc20BalanceByAddress(chain ChainType, contractAddress string, address string) (string, error)
	GetBalanceByAddress(chain ChainType, address string) (string, error)
}

type dappLinkService struct {
	ethClients    map[ChainType]*ethclient.Client
	ethDataClient *EthData
}

func NewDappLinkService(
	conf *config.Config,
	chains ...ChainType,
) (DappLinkService, error) {

	if len(chains) == 0 {
		return nil, errors.New("no chain specified")
	}

	clients := make(map[ChainType]*ethclient.Client)

	for _, chain := range chains {
		rpc, err := rpcByChain(conf, chain)
		if err != nil {
			return nil, err
		}

		cli, err := ethclient.DialContext(context.Background(), rpc)
		if err != nil {
			log.Error("init eth client failed", "chain", chain, "err", err)
			return nil, err
		}

		log.Info("init eth client success", "chain", chain, "rpc", rpc)
		clients[chain] = cli
	}

	ethDataClient, err := NewEthDataClient(
		conf.RpcConfig.DataApiUrl,
		conf.RpcConfig.DataApiKey,
		time.Second*15,
	)
	if err != nil {
		return nil, err
	}

	return &dappLinkService{
		ethClients:    clients,
		ethDataClient: ethDataClient,
	}, nil
}

func (ds *dappLinkService) GetErc20BalanceByAddress(chain ChainType, contractAddress string, address string) (string, error) {

	contractAddr := common.HexToAddress(contractAddress)
	userAddr := common.HexToAddress(address)

	data := BuildErc20BalanceData(userAddr)

	ctxwt, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 构造 CallMsg
	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	res, err := ds.ethClients[chain].CallContract(ctxwt, msg, nil)
	if err != nil {
		log.Error("get balance fail", "err", err)
		return "", err
	}

	if len(res) == 0 {
		return "0", nil
	}

	out := new(big.Int).SetBytes(res)
	fmt.Println(out.String()) // 得到十进制的余额字符串

	return out.String(), nil
}

func (ds *dappLinkService) GetBalanceByAddress(chain ChainType, address string) (string, error) {
	ethClient := ds.ethClients[chain]
	if ethClient == nil {
		return "0", errors.New("eth client not found for chain")
	}
	balanceResult, err := ethClient.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		log.Error("get balance fail", "err", err)
		return "0", err
	}
	log.Info("balance result", "balanceResult", balanceResult.String())
	if balanceResult != nil {
		return balanceResult.String(), nil
	} else {
		return "0", err
	}
}

type EthData struct {
	EthDataCli *etherscan.ChainExplorerAdaptor
}

func NewEthDataClient(baseUrl, apiKey string, timeout time.Duration) (*EthData, error) {
	etherscanCli, err := etherscan.NewChainExplorerAdaptor(apiKey, baseUrl, false, time.Duration(timeout))
	if err != nil {
		log.Error("New etherscan client fail", "err", err)
		return nil, err
	}
	return &EthData{EthDataCli: etherscanCli}, err
}

func (ed *EthData) GetTxByAddress(page, pagesize uint64, address string, action account.ActionType) (*account.TransactionResponse[account.AccountTxResponse], error) {
	request := &account.AccountTxRequest{
		PageRequest: chain.PageRequest{
			Page:  page,
			Limit: pagesize,
		},
		Action:  action,
		Address: address,
	}
	txData, err := ed.EthDataCli.GetTxByAddress(request)
	if err != nil {
		return nil, err
	}
	return txData, nil
}

func (ed *EthData) GetBalanceByAddress(contractAddr, address string) (*account.AccountBalanceResponse, error) {
	accountItem := []string{address}
	symbol := []string{"ETH"}
	contractAddress := []string{contractAddr}
	protocolType := []string{""}
	page := []string{"1"}
	limit := []string{"10"}
	acbr := &account.AccountBalanceRequest{
		ChainShortName:  "ETH",
		ExplorerName:    "etherescan",
		Account:         accountItem,
		Symbol:          symbol,
		ContractAddress: contractAddress,
		ProtocolType:    protocolType,
		Page:            page,
		Limit:           limit,
	}
	etherscanResp, err := ed.EthDataCli.GetAccountBalance(acbr)
	if err != nil {
		log.Error("get account balance error", "err", err)
		return nil, err
	}
	return etherscanResp, nil
}

func BuildErc20BalanceData(address common.Address) []byte {
	var data []byte

	transferFnSignature := []byte("balanceOf(address)")
	hash := crypto.Keccak256Hash(transferFnSignature)
	methodId := hash[:4]
	dataAddress := common.LeftPadBytes(address.Bytes(), 32)

	data = append(data, methodId...)
	data = append(data, dataAddress...)

	return data
}

func rpcByChain(conf *config.Config, chain ChainType) (string, error) {
	switch chain {
	case ChainETH:
		return conf.RpcConfig.EthRpc, nil
	case ChainArbitrum:
		return conf.RpcConfig.ArbitrumRpc, nil
	case ChainPolygon:
		return conf.RpcConfig.PolygonRpc, nil
	case ChainBase:
		return conf.RpcConfig.BaseRpc, nil
	case ChainOp:
		return conf.RpcConfig.OpRpc, nil
	case ChainRootHash:
		return conf.RpcConfig.RootHashRpc, nil
	default:
		return "", fmt.Errorf("unsupported chain: %s", chain)
	}
}
