package entity

import (
	"encoding/json"
	"time"
)

type DepositByWalletID struct {
	WalletID      string `validate:"required" json:"walletID"`
	DepositAmount int    `validate:"required" json:"depositAmount"`
}

type Wallet struct {
	WalletID         string
	Amount           float64
	IsAboveThreshold bool
}

type AboveThreshold struct {
	WalletID      string
	LastCycleTime *time.Time
	Amount        int32
}

func (t AboveThreshold) ToJsonString() string {
	raw, _ := json.Marshal(&t)
	return string(raw)
}

func (t *AboveThreshold) GetAboveThresholdState() bool {
	if t == nil {
		return false
	}

	return t.Amount >= 10000
}
