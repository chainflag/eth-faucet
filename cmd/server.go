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

var (
	apiPortFlag  = flag.Int("apiport", 8080, "Listener port to serve HTTP connection")
	configFlag   = flag.String("config", "config.yml", "Path of wallet config yaml file")
	intervalFlag = flag.Int("interval", 1440, "Number of minutes to wait between funding rounds")
	payoutFlag   = flag.Int("payout", 1, "Number of Ethers to transfer per user request")
	queueCapFlag = flag.Int("queuecap", 100, "Maximum transactions waiting to be sent")
)

func Execute() {
	flag.Parse()
	v := viper.New()
	v.SetConfigFile(*configFlag)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	privateKey, err := getPrivateKey(v.GetStringMapString("wallet"))
	if err != nil {
		panic(fmt.Errorf("failed to read private key: %w", err))
	}

	txBuilder := chain.NewTxBuilder(v.GetString("provider"), privateKey, nil)
	config := server.NewConfig(*apiPortFlag, *intervalFlag, *payoutFlag, *queueCapFlag)
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
