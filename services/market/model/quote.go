package model

import "time"

type Quote struct {
	BaseAsset  string
	QuoteAsset string
	Price      float64
	Volume24h  float64
	Source     string
	Timestamp  time.Time
}
