package server

type Config struct {
	network         string
	symbol          string
	httpPort        int
	interval        int
	payout          float64
	proxyCount      int
	hcaptchaSiteKey string
	hcaptchaSecret  string
	loadBalancing   bool // true if using load balancer
}

func NewConfig(network, symbol string, httpPort, interval, proxyCount int, payout float64, hcaptchaSiteKey, hcaptchaSecret string, loadBalancing bool) *Config {
	return &Config{
		network:         network,
		symbol:          symbol,
		httpPort:        httpPort,
		interval:        interval,
		payout:          payout,
		proxyCount:      proxyCount,
		hcaptchaSiteKey: hcaptchaSiteKey,
		hcaptchaSecret:  hcaptchaSecret,
		loadBalancing:   loadBalancing,
	}
}
