package cmd

import (
	"flag"
	"os"
	"os/signal"

	"github.com/spf13/viper"

	"github.com/chainflag/eth-faucet/internal"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8080, "listen port")
	flag.Parse()
}

func initConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yml")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	return v
}

func Execute() {
	conf := initConfig()
	provider := conf.GetString("provider")
	privKey := conf.GetString("privkey")
	maxQueue := conf.GetInt("maxqueue")

	faucet := internal.NewFaucet(maxQueue, internal.NewTxBuilder(provider, privKey))
	faucet.SetPayoutEther(int64(conf.GetInt("payout")))
	go faucet.Run()

	server := internal.NewServer(faucet)
	go server.Run(port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
