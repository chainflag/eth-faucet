package server

import (
	"math/big"
	"time"
)

type Config struct {
	apiPort    int
	interval   time.Duration
	payout     *big.Int
	proxyCount int
	queueCap   int
}

func NewConfig(apiPort, interval, payout, proxyCount, queueCap int) *Config {
	ether := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	return &Config{
		apiPort:    apiPort,
		interval:   time.Duration(interval),
		payout:     new(big.Int).Mul(big.NewInt(int64(payout)), ether),
		proxyCount: proxyCount,
		queueCap:   queueCap,
	}
}
