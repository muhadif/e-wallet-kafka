package handler

import (
	"context"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"github.com/muhadif/e-wallet/config"
	"github.com/muhadif/e-wallet/core/module"
	"log"
	"sync"
)

func NewProcessor(ctx context.Context, walletModule module.WalletModule, wg *sync.WaitGroup) {
	RunBalanceProcessor(ctx, walletModule, wg)
	RunThresholdProcessor(ctx, walletModule, wg)
}

func RunBalanceProcessor(ctx context.Context, walletModule module.WalletModule, wg *sync.WaitGroup) {
	wg.Add(1)
	balanceGroup := goka.DefineGroup(config.GroupBalance,
		goka.Input(config.TopicDeposit, new(codec.String), walletModule.AddBalanceToWallet),
		goka.Persist(new(codec.String)),
	)

	balanceProcessor, err := goka.NewProcessor(config.Brokers, balanceGroup)
	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}

	go func() {
		defer wg.Done()
		if err = balanceProcessor.Run(ctx); err != nil {
			log.Fatalf("error running processor: %v", err)
		} else {
			log.Printf("Processor shutdown cleanly")
		}
	}()
}

func RunThresholdProcessor(ctx context.Context, walletModule module.WalletModule, wg *sync.WaitGroup) {
	wg.Add(1)
	aboveThresholdGroup := goka.DefineGroup(config.GroupAboveThreshold,
		goka.Input(config.TopicDeposit, new(codec.String), walletModule.AboveThresholdChecker),
		goka.Persist(new(codec.String)),
	)

	aboveThresholdProcessor, err := goka.NewProcessor(config.Brokers, aboveThresholdGroup)
	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}

	go func() {
		defer wg.Done()
		if err = aboveThresholdProcessor.Run(ctx); err != nil {
			log.Fatalf("error running processor: %v", err)
		} else {
			log.Printf("Processor shutdown cleanly")
		}
	}()
}
