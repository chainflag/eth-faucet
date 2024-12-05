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
	logoURL         string
	backgroundURL   string
}

func NewConfig(network, symbol string, httpPort, interval, proxyCount int, payout int64, hcaptchaSiteKey, hcaptchaSecret, logoURL, backgroundURL string) *Config {
	return &Config{
		network:         network,
		symbol:          symbol,
		httpPort:        httpPort,
		interval:        interval,
		payout:          payout,
		proxyCount:      proxyCount,
		hcaptchaSiteKey: hcaptchaSiteKey,
		hcaptchaSecret:  hcaptchaSecret,
		logoURL:         logoURL,
		backgroundURL:   backgroundURL,
	}
}
