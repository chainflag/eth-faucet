package cmd

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/chainflag/eth-faucet/internal/chain"
	"github.com/chainflag/eth-faucet/internal/server"
)

const (
	AppVersion     = "v1.0.0"
	DefaultKeyAuth = "password.txt"
)

var chainIDMap = map[string]int{"mainnet": 1, "ropsten": 3, "rinkeby": 4, "goerli": 5, "kovan": 42}

var (
	chainNameFlag = flag.String("chainname", "testnet", "Network name to display on the frontend")
	httpPortFlag  = flag.Int("httpport", 8080, "Listener port to serve HTTP connection")
	intervalFlag  = flag.Int("interval", 1440, "Number of minutes to wait between funding rounds")
	payoutFlag    = flag.Int("payout", 1, "Number of Ethers to transfer per user request")
	proxyCntFlag  = flag.Int("proxycount", 0, "Count of reverse proxies in front of the server")
	queueCapFlag  = flag.Int("queuecap", 100, "Maximum transactions waiting to be sent")
	versionFlag   = flag.Bool("v", false, "Print version number")
)

func init() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(AppVersion)
		os.Exit(0)
	}
}

func Execute() {
	privateKey, err := getPrivateKeyFromEnv()
	if err != nil {
		panic(fmt.Errorf("failed to read private key: %w", err))
	}
	var chainID *big.Int
	if value, ok := chainIDMap[strings.ToLower(*chainNameFlag)]; ok {
		chainID = big.NewInt(int64(value))
	}

	txBuilder, err := chain.NewTxBuilder(os.Getenv("WEB3_PROVIDER"), privateKey, chainID)
	if err != nil {
		panic(fmt.Errorf("cannot connect to web3 provider: %w", err))
	}
	config := server.NewConfig(*chainNameFlag, *httpPortFlag, *intervalFlag, *payoutFlag, *proxyCntFlag, *queueCapFlag)
	go server.NewServer(txBuilder, config).Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func getPrivateKeyFromEnv() (*ecdsa.PrivateKey, error) {
	if value, ok := os.LookupEnv("PRIVATE_KEY"); ok {
		return crypto.HexToECDSA(value)
	}
	keydir, ok := os.LookupEnv("KEYSTORE")
	if !ok {
		fmt.Println("Please set the environment variable for private key or keystore")
		os.Exit(1)
	}

	keyfile, err := chain.ResolveKeyfilePath(keydir)
	if err != nil {
		panic(err)
	}
	password, err := os.ReadFile(DefaultKeyAuth)
	if err != nil {
		panic(fmt.Errorf("failed to read password from %v", DefaultKeyAuth))
	}

	return chain.DecryptPrivateKey(keyfile, strings.TrimRight(string(password), "\r\n"))
}
