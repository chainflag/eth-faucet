package server

type Config struct {
	network         string
	symbol          string
	httpPort        int
	interval        int
	payout          int64
	proxyCount      int
	hcaptchaSiteKey string
	hcaptchaSecret  string
	logoUrl         string
	backgroundUrl   string
}

func NewConfig(network, symbol string, httpPort, interval, proxyCount int, payout int64, hcaptchaSiteKey, hcaptchaSecret, logoUrl, backgroundUrl string) *Config {
	return &Config{
		network:         network,
		symbol:          symbol,
		httpPort:        httpPort,
		interval:        interval,
		payout:          payout,
		proxyCount:      proxyCount,
		hcaptchaSiteKey: hcaptchaSiteKey,
		hcaptchaSecret:  hcaptchaSecret,
		logoUrl:         logoUrl,
		backgroundUrl:   backgroundUrl,
	}
}
