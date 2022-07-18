package config

import "github.com/lovoo/goka"

const (
	TopicDeposit        goka.Stream = "deposit-wallet"
	GroupBalance        goka.Group  = "balance-wallet"
	GroupAboveThreshold goka.Group  = "above-threshold-wallet"
)

var Brokers = []string{"localhost:29092", "localhost:39092"}
