package cmd

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"

	"github.com/chainflag/eth-faucet/internal"
	"github.com/chainflag/eth-faucet/internal/chain"
)

type config struct {
	Provider   string
	Payout     int
	QueueCap   int
	PrivateKey *ecdsa.PrivateKey
}

var port int

func init() {
	flag.IntVar(&port, "port", 8080, "listen port")
	flag.Parse()
}

func initConfig() *config {
	v := viper.New()
	v.SetConfigFile("./config.yml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %v \n", err))
	}

	privateKey, err := func(walletConf map[string]string) (*ecdsa.PrivateKey, error) {
		if walletConf["privkey"] != "" {
			return crypto.HexToECDSA(walletConf["privkey"])
		}

		keyfile, err := chain.ResolveKeyfilePath(walletConf["keystore"])
		if err != nil {
			panic(err)
		}

		return chain.DecryptPrivateKey(keyfile, walletConf["password"])
	}(v.GetStringMapString("wallet"))

	if err != nil {
		panic(fmt.Errorf("Failed to parse private key: %v \n", err))
	}

	return &config{
		Provider:   v.GetString("provider"),
		Payout:     v.GetInt("payout"),
		QueueCap:   v.GetInt("queuecap"),
		PrivateKey: privateKey,
	}
}

func Execute() {
	conf := initConfig()
	faucet := internal.NewFaucet(chain.NewTxBuilder(conf.Provider, conf.PrivateKey, nil), conf.QueueCap)
	defer faucet.Close()
	faucet.SetPayoutEther(conf.Payout)
	go faucet.Run()

	server := internal.NewServer(faucet)
	go server.Start(port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
