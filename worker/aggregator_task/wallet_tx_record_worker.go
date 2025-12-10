// wallet_tx_record_worker.go
package aggregator_task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"

	dbBackend "github.com/roothash-pay/wallet-services/database/backend"
	"github.com/roothash-pay/wallet-services/services/common/chaininfo"
	"github.com/roothash-pay/wallet-services/services/grpc_client/account"
)

// WalletTxRecordWorkerConfig 配置
type WalletTxRecordWorkerConfig struct {
	// 扫描间隔（秒）
	ScanInterval int
	// 上次检查时间阈值（秒）- 只扫描上次检查时间早于此阈值的记录
	LastCheckedThreshold int
	// 每次扫描的最大记录数
	BatchSize int
	// 并发度（worker pool 大小）
	Concurrency int
	// 超时阈值（秒）- 超过此时间未确认的交易标记为失败
	TimeoutThreshold int
}

// WalletTxRecordWorker 定时扫描 pending 交易并更新状态
type WalletTxRecordWorker struct {
	db            dbBackend.WalletTxRecordDB
	accountClient *account.WalletAccountClient
	chainInfo     chaininfo.Provider
	config        WalletTxRecordWorkerConfig
	stopCh        chan struct{}
	wg            sync.WaitGroup
}

// NewWalletTxRecordWorker 创建 worker
func NewWalletTxRecordWorker(
	db dbBackend.WalletTxRecordDB,
	accountClient *account.WalletAccountClient,
	chainInfo chaininfo.Provider,
	config WalletTxRecordWorkerConfig,
) *WalletTxRecordWorker {
	// 设置默认值
	if config.ScanInterval <= 0 {
		config.ScanInterval = 10 // 默认 10 秒
	}
	if config.LastCheckedThreshold <= 0 {
		config.LastCheckedThreshold = 5 // 默认 5 秒
	}
	if config.BatchSize <= 0 {
		config.BatchSize = 100 // 默认 100 条
	}
	if config.Concurrency <= 0 {
		config.Concurrency = 10 // 默认 10 个并发
	}
	if config.TimeoutThreshold <= 0 {
		config.TimeoutThreshold = 3600 // 默认 1 小时
	}

	return &WalletTxRecordWorker{
		db:            db,
		accountClient: accountClient,
		chainInfo:     chainInfo,
		config:        config,
		stopCh:        make(chan struct{}),
	}
}

// Start 启动 worker
func (w *WalletTxRecordWorker) Start() {
	w.wg.Add(1)
	go w.run()
	log.Info("WalletTxRecordWorker started",
		"scanInterval", w.config.ScanInterval,
		"concurrency", w.config.Concurrency,
		"batchSize", w.config.BatchSize)
}

// Stop 停止 worker
func (w *WalletTxRecordWorker) Stop() {
	close(w.stopCh)
	w.wg.Wait()
	log.Info("WalletTxRecordWorker stopped")
}

// run 主循环
func (w *WalletTxRecordWorker) run() {
	defer w.wg.Done()

	ticker := time.NewTicker(time.Duration(w.config.ScanInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.stopCh:
			return
		case <-ticker.C:
			w.scanAndUpdate()
		}
	}
}

// scanAndUpdate 扫描并更新 pending 交易
func (w *WalletTxRecordWorker) scanAndUpdate() {
	ctx := context.Background()

	// 计算上次检查时间阈值
	lastCheckedBefore := time.Now().Add(-time.Duration(w.config.LastCheckedThreshold) * time.Second)

	// 查询需要检查的 pending 交易
	records, err := w.db.GetPendingTxsForCheck(lastCheckedBefore, w.config.BatchSize)
	if err != nil {
		log.Error("Failed to get pending txs for check", "err", err)
		return
	}

	if len(records) == 0 {
		return
	}

	log.Info("Found pending txs to check", "count", len(records))

	// 使用 worker pool 并发处理
	jobs := make(chan *dbBackend.WalletTxRecord, len(records))
	results := make(chan struct{}, len(records))

	// 启动 workers
	for i := 0; i < w.config.Concurrency; i++ {
		w.wg.Add(1)
		go w.worker(ctx, jobs, results)
	}

	// 发送任务
	for _, record := range records {
		jobs <- record
	}
	close(jobs)

	// 等待所有任务完成
	for i := 0; i < len(records); i++ {
		<-results
	}

	log.Info("Finished checking pending txs", "count", len(records))
}

// worker 处理单个交易
func (w *WalletTxRecordWorker) worker(ctx context.Context, jobs <-chan *dbBackend.WalletTxRecord, results chan<- struct{}) {
	defer w.wg.Done()

	for record := range jobs {
		w.checkAndUpdateTx(ctx, record)
		results <- struct{}{}
	}
}

// checkAndUpdateTx 检查并更新单个交易状态
func (w *WalletTxRecordWorker) checkAndUpdateTx(ctx context.Context, record *dbBackend.WalletTxRecord) {
	// 更新 last_checked_at
	now := time.Now()
	defer func() {
		updates := map[string]interface{}{
			"last_checked_at": now,
		}
		_ = w.db.UpdateWalletTxRecord(record.Guid, updates)
	}()

	// 检查是否超时
	if w.isTimeout(record) {
		w.markAsFailed(record, dbBackend.FailReasonNotFoundTimeout, "Transaction not found and timeout")
		return
	}

	info, err := w.getChainInfo(ctx, record.ChainID)
	if err != nil {
		log.Warn("Chain info not available for pending tx", "guid", record.Guid, "chainID", record.ChainID, "err", err)
		return
	}

	txInfo, err := w.accountClient.GetTxByHash(
		ctx,
		info.ConsumerToken,
		info.WalletChain,
		info.WalletCoin,
		info.WalletNetwork,
		record.Hash,
	)
	if err != nil {
		log.Warn("Failed to get tx by hash", "guid", record.Guid, "hash", record.Hash, "err", err)
		return
	}

	if txInfo == nil {
		log.Warn("Tx not found", "guid", record.Guid, "hash", record.Hash)
		return
	}

	// 根据链上状态更新数据库
	// TxStatus: 0=NotFound, 1=Pending, 2=Failed, 3=Success, 4=ContractExecuteFailed
	if txInfo.Status == 3 { // pb.TxStatus_Success
		// 链上确认成功
		w.markAsSuccess(record, txInfo.Height, txInfo.Datetime)
	} else if txInfo.Status == 2 || txInfo.Status == 4 { // pb.TxStatus_Failed or ContractExecuteFailed
		// 链上执行失败
		w.markAsFailed(record, dbBackend.FailReasonChainFailed, "Transaction failed on chain")
	}
}

// isTimeout 检查交易是否超时
func (w *WalletTxRecordWorker) isTimeout(record *dbBackend.WalletTxRecord) bool {
	// 从 created_at 开始计算
	timeout := time.Duration(w.config.TimeoutThreshold) * time.Second
	return time.Since(record.CreateTime) > timeout
}

// markAsSuccess 标记交易为成功
func (w *WalletTxRecordWorker) markAsSuccess(record *dbBackend.WalletTxRecord, blockHeight string, txTime string) {
	updates := map[string]interface{}{
		"status":       dbBackend.TxStatusSuccess,
		"block_height": blockHeight,
		"tx_time":      txTime,
		"memo":         w.updateMemoStatus(record.Memo, "Success"),
	}

	if err := w.db.UpdateWalletTxRecord(record.Guid, updates); err != nil {
		log.Error("Failed to mark tx as success", "guid", record.Guid, "hash", record.Hash, "err", err)
	} else {
		log.Info("Tx marked as success", "guid", record.Guid, "hash", record.Hash, "blockHeight", blockHeight)
	}
}

// markAsFailed 标记交易为失败
func (w *WalletTxRecordWorker) markAsFailed(record *dbBackend.WalletTxRecord, failReasonCode string, failReasonMsg string) {
	updates := map[string]interface{}{
		"status":           dbBackend.TxStatusFailed,
		"fail_reason_code": failReasonCode,
		"fail_reason_msg":  failReasonMsg,
		"memo":             w.updateMemoStatus(record.Memo, "Failed: "+failReasonMsg),
	}

	if err := w.db.UpdateWalletTxRecord(record.Guid, updates); err != nil {
		log.Error("Failed to mark tx as failed", "guid", record.Guid, "hash", record.Hash, "err", err)
	} else {
		log.Info("Tx marked as failed", "guid", record.Guid, "hash", record.Hash, "reason", failReasonCode)
	}
}

// updateMemoStatus 更新 memo 中的状态
func (w *WalletTxRecordWorker) updateMemoStatus(memo string, newStatus string) string {
	// 简单替换：将 (Pending) 或 (Created) 替换为新状态
	// TODO: 更精确的替换逻辑
	if len(memo) > 0 {
		// 移除旧状态
		memo = memo[:len(memo)-len(" (Pending)")]
		if len(memo) > len(" (Created)") && memo[len(memo)-len(" (Created)"):] == " (Created)" {
			memo = memo[:len(memo)-len(" (Created)")]
		}
		// 添加新状态
		return memo + " (" + newStatus + ")"
	}
	return memo
}

func (w *WalletTxRecordWorker) getChainInfo(ctx context.Context, chainID string) (*chaininfo.Info, error) {
	if w.chainInfo == nil {
		return nil, fmt.Errorf("chain info provider not configured")
	}
	return w.chainInfo.Get(ctx, chainID)
}
