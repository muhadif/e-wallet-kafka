package module

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"github.com/lovoo/goka"
	"github.com/muhadif/e-wallet/core/entity"
	walletContract "github.com/muhadif/e-wallet/proto"
	"log"
	"time"
)

type WalletModule interface {
	DepositBalanceByWalletID(ctx context.Context, req *entity.DepositByWalletID) error
	GetBalanceByWalletID(ctx context.Context, walletID string) (*entity.Wallet, error)

	AddBalanceToWallet(ctx goka.Context, msg interface{})
	AboveThresholdChecker(ctx goka.Context, msg interface{})
}

func NewBalanceModule(
	emitter *goka.Emitter,
	balanceView *goka.View,
	aboveThresholdView *goka.View) WalletModule {
	return &wallet{emitter: emitter, balanceView: balanceView, aboveThresholdView: aboveThresholdView}
}

type wallet struct {
	emitter            *goka.Emitter
	balanceView        *goka.View
	aboveThresholdView *goka.View
}

func (w *wallet) DepositBalanceByWalletID(ctx context.Context, req *entity.DepositByWalletID) error {
	depositPayload := &walletContract.Deposit{
		Amount:   int32(req.DepositAmount),
		WalletID: req.WalletID,
	}

	raw, err := json2.Marshal(depositPayload)
	if err != nil {
		return err
	}

	if err := w.emitter.EmitSync(req.WalletID, string(raw)); err != nil {
		return err
	}
	return nil
}

func (w *wallet) GetBalanceByWalletID(ctx context.Context, walletID string) (*entity.Wallet, error) {
	balanceRaw, err := w.balanceView.Get(walletID)
	if err != nil {
		return nil, err
	}

	balance := new(walletContract.Balance)
	if balanceRaw != nil {
		err := json2.Unmarshal([]byte(balanceRaw.(string)), &balance)
		if err != nil {
			log.Fatalf("error unmarsal balance %v", err.Error())
		}
	}

	aboveThresholdRaw, err := w.aboveThresholdView.Get(walletID)
	if err != nil {
		return nil, err
	}

	aboveThresholdData := new(walletContract.AboveThreshold)
	if aboveThresholdRaw != nil {
		err = json2.Unmarshal([]byte(aboveThresholdRaw.(string)), &aboveThresholdData)
		if err != nil {
			log.Fatalf("error unmarsal aboveThresholdData %v", err.Error())
		}
	}

	var isAboveThreshold bool
	if aboveThresholdData.Amount >= 10000 {
		isAboveThreshold = true
	}

	return &entity.Wallet{
		WalletID:         walletID,
		Amount:           float64(balance.Amount),
		IsAboveThreshold: isAboveThreshold,
	}, nil
}

// AddBalanceToWallet add deposit to balance group
func (w *wallet) AddBalanceToWallet(ctx goka.Context, msg interface{}) {
	var deposit *walletContract.Deposit
	_ = json2.Unmarshal([]byte(msg.(string)), &deposit)

	var balance *walletContract.Balance
	if val := ctx.Value(); val != nil {
		_ = json2.Unmarshal([]byte(ctx.Value().(string)), &balance)
		balance.Amount += deposit.Amount
	} else {
		balance = &walletContract.Balance{
			Amount:   deposit.Amount,
			WalletID: deposit.WalletID,
		}
	}
	raw, _ := json2.Marshal(balance)
	ctx.SetValue(string(raw))
}

// AboveThresholdChecker add above threshold chcker and insest to above threshold group
func (w *wallet) AboveThresholdChecker(ctx goka.Context, msg interface{}) {
	currentTime := ctx.Timestamp()

	var deposit *walletContract.Deposit
	_ = json2.Unmarshal([]byte(msg.(string)), &deposit)

	var aboveThresholdData *walletContract.AboveThreshold
	if val := ctx.Value(); val != nil {
		_ = json2.Unmarshal([]byte(val.(string)), &aboveThresholdData)
	} else {
		aboveThresholdData = &walletContract.AboveThreshold{
			Amount:    deposit.Amount,
			WalletID:  deposit.WalletID,
			CycleTime: currentTime.Format(time.RFC3339),
		}
	}

	lastCycleTime, _ := time.Parse(time.RFC3339, aboveThresholdData.CycleTime)
	switch {
	case lastCycleTime.Unix() == currentTime.Unix():
		aboveThresholdData.Amount += deposit.Amount
	case lastCycleTime.Add(2 * time.Minute).After(currentTime):
		aboveThresholdData.Amount += deposit.Amount
	case lastCycleTime.Add(2 * time.Minute).Before(currentTime):
		aboveThresholdData.CycleTime = currentTime.Format(time.RFC3339)
		aboveThresholdData.Amount = deposit.Amount
	}

	log.Println(fmt.Sprintf("%v", aboveThresholdData))
	raw, _ := json2.Marshal(aboveThresholdData)
	ctx.SetValue(string(raw))
}
