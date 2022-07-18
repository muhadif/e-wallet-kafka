package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"github.com/muhadif/e-wallet/config"
	"github.com/muhadif/e-wallet/core/module"
	"github.com/muhadif/e-wallet/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	err := checkKafkaTopics()
	if err != nil {
		panic(fmt.Sprintf("error check kafka topic : %s", err.Error()))
	}

	balanceEmitter, err := goka.NewEmitter(config.Brokers, config.TopicDeposit, new(codec.String))
	if err != nil {
		panic(fmt.Sprintf("error init broker %s", err.Error()))
	}

	balanceView, err := goka.NewView(config.Brokers,
		goka.GroupTable(config.GroupBalance),
		new(codec.String),
	)

	thresholdView, err := goka.NewView(config.Brokers,
		goka.GroupTable(config.GroupAboveThreshold),
		new(codec.String),
	)
	if err != nil {
		panic(fmt.Sprintf("error init view %s", err.Error()))
	}

	walletModule := module.NewBalanceModule(balanceEmitter, balanceView, thresholdView)
	walletHandler := handler.NewWalletHandlerAPI(walletModule)

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	handler.NewProcessor(ctx, walletModule, &wg)

	router := mux.NewRouter()
	router.HandleFunc("/api/deposit", walletHandler.DepositWalletByWalletID).Methods("POST")
	router.HandleFunc("/api/wallet/{walletID}", walletHandler.GetDepositByWalletID).Methods("GET")

	go balanceView.Run(ctx)
	go thresholdView.Run(ctx)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
	}

	go srv.ListenAndServe()

	log.Println("All component is running OK, enjoy!")

	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
	<-wait

	cancel()

	wg.Wait()
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

}

func checkKafkaTopics() error {
	tmc := goka.NewTopicManagerConfig()
	tmc.Table.Replication = 1
	tmc.Stream.Replication = 1

	tm, err := goka.NewTopicManager(config.Brokers, goka.DefaultConfig(), tmc)
	if err != nil {
		return err
	}

	err = tm.EnsureStreamExists(string(config.TopicDeposit), 1)
	if err != nil {
		return err
	}

	err = tm.EnsureTableExists(fmt.Sprintf("%s-table", string(config.GroupBalance)), 1)
	if err != nil {
		return err
	}

	err = tm.EnsureTableExists(fmt.Sprintf("%s-table", string(config.GroupAboveThreshold)), 1)
	if err != nil {
		return err
	}

	return nil
}
