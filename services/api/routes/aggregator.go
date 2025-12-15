package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/models/backend"
	"github.com/roothash-pay/wallet-services/services/api/service"
)

// AggregatorRoutes handles swap aggregator related routes
type AggregatorRoutes struct {
	aggregatorService *service.AggregatorService
}

// NewAggregatorRoutes creates a new aggregator routes handler
func NewAggregatorRoutes(aggregatorService *service.AggregatorService) *AggregatorRoutes {
	return &AggregatorRoutes{
		aggregatorService: aggregatorService,
	}
}

// RegisterRoutes registers all aggregator routes to the router
func (h *AggregatorRoutes) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/aggregator", func(r chi.Router) {
		r.Post("/quotes", h.GetQuotesHandler)
		r.Post("/swap/prepare", h.PrepareSwapHandler)
		r.Post("/tx/submitSigned", h.SubmitSignedTxHandler)
		r.Get("/swap/status", h.GetSwapStatusHandler)
	})
}

// GetQuotesHandler godoc
// @Summary      获取聚合报价
// @Description  并发调用所有可用 provider，返回最佳报价及候选列表
// @Tags         Aggregator
// @Accept       json
// @Produce      json
// @Param        request  body      backend.QuoteRequest true "报价请求"
// @Success      200      {object}  backend.QuoteResponse
// @Failure      400      {string}  string "invalid request body"
// @Failure      500      {string}  string "internal error"
// @Router       /aggregator/quotes [post]
func (h *AggregatorRoutes) GetQuotesHandler(w http.ResponseWriter, r *http.Request) {
	var req backend.QuoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.aggregatorService.GetQuotes(r.Context(), &req)
	if err != nil {
		log.Error("GetQuotes failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// PrepareSwapHandler godoc
// @Summary      生成 swap 执行计划
// @Description  根据报价 ID 和用户地址生成签名动作链路
// @Tags         Aggregator
// @Accept       json
// @Produce      json
// @Param        request  body      backend.PrepareSwapRequest true "prepare 请求"
// @Success      200      {object}  backend.PrepareSwapResponse
// @Failure      400      {string}  string "invalid request body"
// @Failure      500      {string}  string "internal error"
// @Router       /aggregator/swap/prepare [post]
func (h *AggregatorRoutes) PrepareSwapHandler(w http.ResponseWriter, r *http.Request) {
	var req backend.PrepareSwapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	quoteID := req.QuoteID
	bestQuotesIndex := req.BestQuotesIndex

	resp, err := h.aggregatorService.PrepareSwap(r.Context(), quoteID, bestQuotesIndex)
	if err != nil {
		log.Error("PrepareSwap failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// SubmitSignedTxHandler godoc
// @Summary      提交签名交易
// @Description  接收客户端签名后的原始交易并广播，同时更新步骤状态
// @Tags         Aggregator
// @Accept       json
// @Produce      json
// @Param        request  body      backend.SubmitSignedTxRequest true "签名交易请求"
// @Success      200      {object}  backend.SubmitSignedTxResponse
// @Failure      400      {string}  string "invalid request body"
// @Failure      500      {string}  string "internal error"
// @Router       /aggregator/tx/submitSigned [post]
func (h *AggregatorRoutes) SubmitSignedTxHandler(w http.ResponseWriter, r *http.Request) {
	var req backend.SubmitSignedTxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.aggregatorService.SubmitSignedTx(r.Context(), &req)
	if err != nil {
		log.Error("SubmitSignedTx failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// SubmitTxHashHandler godoc
// @Summary      提交交易哈希（前端钱包已广播）
// @Description  用户使用前端钱包（如 MetaMask）自行广播交易后，将 txHash 回传给后端。后端记录 swap step 的 txHash 并将状态置为 PENDING，后续可通过 GetSwapStatus/worker 跟踪链上结果。
// @Tags         Aggregator
// @Accept       json
// @Produce      json
// @Param        request  body      backend.SubmitTxHashRequest  true  "提交交易哈希请求"
// @Success      200      {object}  backend.SubmitTxHashResponse
// @Failure      400      {string}  string "invalid request body"
// @Failure      500      {string}  string "internal error"
// @Router       /aggregator/swap/submit_tx_hash [post]
func (h *AggregatorRoutes) SubmitTxHashHandler(w http.ResponseWriter, r *http.Request) {
	var req backend.SubmitTxHashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.aggregatorService.SubmitTxHash(r.Context(), &req)
	if err != nil {
		log.Error("SubmitTxHash failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetSwapStatusHandler godoc
// @Summary      查询 swap 状态
// @Description  根据 swapId 返回各步骤最新状态，并在需要时刷新链上信息
// @Tags         Aggregator
// @Produce      json
// @Param        swapId  query     string true "Swap ID"
// @Success      200     {object}  backend.SwapStatusResponse
// @Failure      400     {string}  string "swapId is required"
// @Failure      500     {string}  string "internal error"
// @Router       /aggregator/swap/status [get]
func (h *AggregatorRoutes) GetSwapStatusHandler(w http.ResponseWriter, r *http.Request) {
	swapID := r.URL.Query().Get("swapId")
	if swapID == "" {
		http.Error(w, "swapId is required", http.StatusBadRequest)
		return
	}

	resp, err := h.aggregatorService.GetSwapStatus(r.Context(), swapID)
	if err != nil {
		log.Error("GetSwapStatus failed", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
