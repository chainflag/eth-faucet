package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/chainflag/eth-faucet/core"
	"github.com/chainflag/eth-faucet/util/conf"
)

func init() {
	flag.IntVar(&port, "port", 8080, "listen port")
	flag.Parse()
}

func main() {
	provider := conf.GetString("provider")
	privKey := conf.GetString("privkey")

	faucet := core.NewFaucet(core.NewTxBuilder(provider, privKey))
	faucet.SetPayoutEther(int64(conf.GetInt("payout")))

	go NewServer(faucet).Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
