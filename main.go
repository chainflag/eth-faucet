package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/chainflag/eth-faucet/core"
	"github.com/chainflag/eth-faucet/util/conf"
)

var (
	port      int
	queueSize int
)

func init() {
	flag.IntVar(&port, "port", 8080, "listen port")
	flag.IntVar(&queueSize, "cap", 100, "queue size")
	flag.Parse()
}

func main() {
	provider := conf.GetString("provider")
	privKey := conf.GetString("privkey")

	faucet := core.NewFaucet(core.NewTxBuilder(provider, privKey))
	faucet.SetPayoutEther(int64(conf.GetInt("payout")))

	worker := NewWorker(queueSize, faucet)
	server := NewServer(worker)
	go worker.Run()
	go server.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
