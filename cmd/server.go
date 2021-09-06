package cmd

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"

	"github.com/chainflag/eth-faucet/internal/chain"
	"github.com/chainflag/eth-faucet/internal/server"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8080, "listen port")
	flag.Parse()
}

func Execute() {
	v := viper.New()
	v.SetConfigFile("config.yml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	privateKey, err := getPrivateKey(v.GetStringMapString("wallet"))
	if err != nil {
		panic(fmt.Errorf("failed to read private key: %w", err))
	}

	txBuilder := chain.NewTxBuilder(v.GetString("provider"), privateKey, nil)
	config := server.NewConfig(v.GetInt("payout"), v.GetInt("queuecap"), port)
	go server.NewServer(txBuilder, config).Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func getPrivateKey(walletConf map[string]string) (*ecdsa.PrivateKey, error) {
	if walletConf["privkey"] != "" {
		return crypto.HexToECDSA(walletConf["privkey"])
	}
	keyfile, err := chain.ResolveKeyfilePath(walletConf["keystore"])
	if err != nil {
		panic(err)
	}
	return chain.DecryptPrivateKey(keyfile, walletConf["password"])
}
