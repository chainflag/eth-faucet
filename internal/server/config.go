package server

import (
	"math/big"
)

type Config struct {
	network           string
	symbol            string
	httpPort          int
	interval          int
	payout            int64
	proxyCount        int
	hcaptchaSiteKey   string
	hcaptchaSecret    string
	logoURL           string
	backgroundURL     string
	frontendType      string
	paidCustomer      bool
	mainnetProvider   string
	minMainnetBalance *big.Int
}

func NewConfig(
	network, symbol string,
	httpPort, interval, proxyCount int,
	payout int64,
	hcaptchaSiteKey, hcaptchaSecret, logoURL, backgroundURL string,
	frontendType string,
	paidCustomer bool,
	mainnetProvider string,
	minMainnetBalance *big.Int,
) *Config {
	return &Config{
		network:           network,
		symbol:            symbol,
		httpPort:          httpPort,
		interval:          interval,
		payout:            payout,
		proxyCount:        proxyCount,
		hcaptchaSiteKey:   hcaptchaSiteKey,
		hcaptchaSecret:    hcaptchaSecret,
		logoURL:           logoURL,
		backgroundURL:     backgroundURL,
		frontendType:      frontendType,
		paidCustomer:      paidCustomer,
		mainnetProvider:   mainnetProvider,
		minMainnetBalance: minMainnetBalance,
	}
}
