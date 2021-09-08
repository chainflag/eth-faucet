package cmd

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/yaml.v2"

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

func init() {
	flag.Parse()
}

func Execute() {
	wallet := &Wallet{}
	if err := loadWallet(*configFlag, wallet); err != nil {
		panic(err)
	}
	privateKey, err := wallet.readPrivateKey()
	if err != nil {
		panic(fmt.Errorf("failed to read private key: %w", err))
	}

	txBuilder := chain.NewTxBuilder(wallet.Provider, privateKey, nil)
	config := server.NewConfig(*apiPortFlag, *intervalFlag, *payoutFlag, *queueCapFlag)
	go server.NewServer(txBuilder, config).Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

type Wallet struct {
	Provider string
	Wallet   struct {
		PrivateKey string `yaml:"privkey"`
		Keystore   string `yaml:"keystore"`
		Password   string `yaml:"password"`
	}
}

func loadWallet(path string, wallet *Wallet) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return yaml.NewDecoder(file).Decode(wallet)
}

func (w Wallet) readPrivateKey() (*ecdsa.PrivateKey, error) {
	wallet := w.Wallet
	if wallet.PrivateKey != "" {
		return crypto.HexToECDSA(wallet.PrivateKey)
	}
	keyfile, err := chain.ResolveKeyfilePath(wallet.Keystore)
	if err != nil {
		panic(err)
	}
	return chain.DecryptPrivateKey(keyfile, wallet.Password)
}
