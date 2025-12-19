package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/roothash-pay/wallet-services/services/market/model"
)

type BinanceProvider struct {
	client *http.Client
}

func NewBinanceProvider() *BinanceProvider {
	return &BinanceProvider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *BinanceProvider) Name() string {
	return "binance"
}

type binanceTicker struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"lastPrice"`
	Volume    string `json:"volume"`
	CloseTime int64  `json:"closeTime"`
}

func (p *BinanceProvider) FetchQuotes(ctx context.Context) ([]model.Quote, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://api.binance.com/api/v3/ticker/24hr",
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 非 200，直接失败
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"binance http %d, body=%s",
			resp.StatusCode,
			string(body[:min(len(body), 200)]),
		)
	}

	// Content-Type 必须是 JSON
	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"binance invalid content-type=%s body=%s",
			ct,
			string(body[:min(len(body), 200)]),
		)
	}

	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	// 尝试解析成数组（正常路径）
	var tickers []binanceTicker
	if err := json.Unmarshal(raw, &tickers); err != nil {
		// 再尝试解析成错误对象（兜底路径）
		var apiErr struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}
		if err2 := json.Unmarshal(raw, &apiErr); err2 == nil && apiErr.Code != 0 {
			return nil, fmt.Errorf("binance api error: %d %s", apiErr.Code, apiErr.Msg)
		}
		// 既不是数组，也不是标准错误，直接忽略
		return nil, err
	}

	quotes := make([]model.Quote, 0, len(tickers))

	for _, t := range tickers {
		// 只要 USDT 计价的主流币（先简单）
		if !strings.HasSuffix(t.Symbol, "USDT") {
			continue
		}

		price, err1 := strconv.ParseFloat(t.LastPrice, 64)
		volume, err2 := strconv.ParseFloat(t.Volume, 64)
		if err1 != nil || err2 != nil {
			continue
		}

		quotes = append(quotes, model.Quote{
			BaseAsset:  strings.TrimSuffix(t.Symbol, "USDT"),
			QuoteAsset: "USDT",
			Price:      price,
			Volume24h:  volume,
			Source:     p.Name(),
			Timestamp:  time.UnixMilli(t.CloseTime),
		})
	}

	return quotes, nil
}
