package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/roothash-pay/wallet-services/services/market/model"
)

type OkxProvider struct {
	client *http.Client
}

func NewOkxProvider() *OkxProvider {
	return &OkxProvider{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *OkxProvider) Name() string {
	return "okx"
}

/*
OKX 返回结构示例：

{
  "code": "0",
  "msg": "",
  "data": [
    {
      "instId": "BTC-USDT",
      "last": "43123.5",
      "vol24h": "123456.78",
      "ts": "1700000000000"
    }
  ]
}
*/

type okxTickerResp struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data []okxTicker `json:"data"`
}

type okxTicker struct {
	InstID string `json:"instId"`
	Last   string `json:"last"`
	Vol24h string `json:"vol24h"`
	TS     string `json:"ts"`
}

func (p *OkxProvider) FetchQuotes(ctx context.Context) ([]model.Quote, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://www.okx.com/api/v5/market/tickers?instType=SPOT",
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

	var result okxTickerResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// OKX code != 0 表示失败
	if result.Code != "0" {
		return nil, nil // 信息级行情：失败直接跳过
	}

	quotes := make([]model.Quote, 0, 16)

	for _, t := range result.Data {
		// 只处理 USDT 交易对
		if !strings.HasSuffix(t.InstID, "-USDT") {
			continue
		}

		price, err1 := strconv.ParseFloat(t.Last, 64)
		volume, err2 := strconv.ParseFloat(t.Vol24h, 64)
		ts, err3 := strconv.ParseInt(t.TS, 10, 64)

		if err1 != nil || err2 != nil || err3 != nil {
			continue
		}

		base := strings.TrimSuffix(t.InstID, "-USDT")

		quotes = append(quotes, model.Quote{
			BaseAsset:  base,
			QuoteAsset: "USDT",
			Price:      price,
			Volume24h:  volume,
			Source:     p.Name(),
			Timestamp:  time.UnixMilli(ts),
		})
	}

	return quotes, nil
}
