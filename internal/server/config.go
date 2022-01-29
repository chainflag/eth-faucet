package server

type Config struct {
	chainName  string
	httpPort   int
	interval   int
	payout     int
	proxyCount int
	queueCap   int
}

func NewConfig(chainName string, httpPort, interval, payout, proxyCount, queueCap int) *Config {
	return &Config{
		chainName:  chainName,
		httpPort:   httpPort,
		interval:   interval,
		payout:     payout,
		proxyCount: proxyCount,
		queueCap:   queueCap,
	}
}
