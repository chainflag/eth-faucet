package server

type Config struct {
	network    string
	symbol     string
	httpPort   int
	interval   int
	payout     int
	proxyCount int
	queueCap   int
}

func NewConfig(network, symbol string, httpPort, interval, payout, proxyCount, queueCap int) *Config {
	return &Config{
		network:    network,
		symbol:     symbol,
		httpPort:   httpPort,
		interval:   interval,
		payout:     payout,
		proxyCount: proxyCount,
		queueCap:   queueCap,
	}
}
