package server

type Config struct {
	network    string
	httpPort   int
	interval   int
	payout     int
	proxyCount int
	queueCap   int
}

func NewConfig(network string, httpPort, interval, payout, proxyCount, queueCap int) *Config {
	return &Config{
		network:    network,
		httpPort:   httpPort,
		interval:   interval,
		payout:     payout,
		proxyCount: proxyCount,
		queueCap:   queueCap,
	}
}
