package server

import (
	"math/big"
)

type Config struct {
	payout   *big.Int
	queueCap int
	port     int
}

func NewConfig(payout, queueCap, port int) *Config {
	ether := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	return &Config{
		payout:   new(big.Int).Mul(big.NewInt(int64(payout)), ether),
		queueCap: queueCap,
		port:     port,
	}
}
