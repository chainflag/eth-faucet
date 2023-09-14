package server

type Config struct {
	network         string
	symbol          string
	httpPort        int
	interval        int
	payout          int
	proxyCount      int
	queueCap        int
	hcaptchaSiteKey string
	hcaptchaSecret  string
}

func NewConfig(network, symbol string, httpPort, interval, payout, proxyCount, queueCap int, hcaptchaSiteKey, hcaptchaSecret string) *Config {
	return &Config{
		network:         network,
		symbol:          symbol,
		httpPort:        httpPort,
		interval:        interval,
		payout:          payout,
		proxyCount:      proxyCount,
		queueCap:        queueCap,
		hcaptchaSiteKey: hcaptchaSiteKey,
		hcaptchaSecret:  hcaptchaSecret,
	}
}
