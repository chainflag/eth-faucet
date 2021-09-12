package cmd

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/chainflag/eth-faucet/internal/chain"
	"github.com/chainflag/eth-faucet/internal/server"
)

var (
	apiPortFlag  = flag.Int("apiport", 8080, "Listener port to serve HTTP connection")
	intervalFlag = flag.Int("interval", 1440, "Number of minutes to wait between funding rounds")
	payoutFlag   = flag.Int("payout", 1, "Number of Ethers to transfer per user request")
	proxyCntFlag = flag.Int("proxycount", 0, "Count of reverse proxies in front of the server")
	queueCapFlag = flag.Int("queuecap", 100, "Maximum transactions waiting to be sent")
)

func init() {
	flag.Parse()
}

func Execute() {
	privateKey, err := getPrivateKeyFromEnv()
	if err != nil {
		panic(fmt.Errorf("failed to read private key: %w", err))
	}

	txBuilder := chain.NewTxBuilder(os.Getenv("WEB3_PROVIDER"), privateKey, nil)
	config := server.NewConfig(*apiPortFlag, *intervalFlag, *payoutFlag, *proxyCntFlag, *queueCapFlag)
	go server.NewServer(txBuilder, config).Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func getPrivateKeyFromEnv() (*ecdsa.PrivateKey, error) {
	if os.Getenv("PRIVATE_KEY") != "" {
		return crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	}
	keyfile, err := chain.ResolveKeyfilePath(os.Getenv("KEYSTORE"))
	if err != nil {
		panic(err)
	}
	return chain.DecryptPrivateKey(keyfile, os.Getenv("PASSWORD"))
}
