package cmd

import (
	"flag"
	"os"
	"os/signal"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"

	"github.com/chainflag/eth-faucet/internal"
	"github.com/chainflag/eth-faucet/internal/chain"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8080, "listen port")
	flag.Parse()
}

func initConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigFile("./config.yml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	return v
}

func Execute() {
	conf := initConfig()
	provider := conf.GetString("provider")
	queueCap := conf.GetInt("queuecap")

	privateKey, err := crypto.HexToECDSA(conf.GetString("privkey"))
	if err != nil {
		panic(err)
	}

	faucet := internal.NewFaucet(chain.NewTxBuilder(provider, privateKey, nil), queueCap)
	defer faucet.Close()
	faucet.SetPayoutEther(int64(conf.GetInt("payout")))
	go faucet.Run()

	server := internal.NewServer(faucet)
	go server.Start(port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
