package cmd

import (
	"crypto/ecdsa"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/chainflag/eth-faucet/internal/chain"
	"github.com/chainflag/eth-faucet/internal/server"
)

var (
	appVersion = "v1.1.0"
	chainIDMap = map[string]int{"goerli": 5, "auroria": 205205, "sepolia": 11155111}

	//httpPortFlag = flag.Int("httpport", 8080, "Listener port to serve HTTP connection")
	proxyCntFlag = flag.Int("proxycount", 0, "Count of reverse proxies in front of the server")
	versionFlag  = flag.Bool("version", false, "Print version number")

	//payoutFlag   = flag.Int("faucet_amount", 10000, "Number of Ethers to transfer per user request")
	intervalFlag = flag.Int("faucet_minutes", 1440, "Number of minutes to wait between funding rounds")
	netnameFlag  = flag.String("faucet_name", os.Getenv("FAUCET_NAME"), "Network name to display on the frontend")
	symbolFlag   = flag.String("faucet_symbol", os.Getenv("FAUCET_SYMBOL"), "Token symbol to display on the frontend")

	keyJSONFlag  = flag.String("wallet_keyjson", os.Getenv("KEYSTORE"), "Keystore file to fund user requests with")
	keyPassFlag  = flag.String("wallet_keypass", "password.txt", "Passphrase text file to decrypt keystore")
	privKeyFlag  = flag.String("wallet_privkey", os.Getenv("PRIVATE_KEY"), "Private key hex to fund user requests with")
	providerFlag = flag.String("wallet_provider", os.Getenv("WEB3_PROVIDER"), "Endpoint for Ethereum JSON-RPC connection")

	hcaptchaSiteKeyFlag = flag.String("hcaptcha_sitekey", os.Getenv("HCAPTCHA_SITEKEY"), "hCaptcha sitekey")
	hcaptchaSecretFlag  = flag.String("hcaptcha_secret", os.Getenv("HCAPTCHA_SECRET"), "hCaptcha secret")
)

func init() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(appVersion)
		os.Exit(0)
	}
}

func Execute() {
	privateKey, err := getPrivateKeyFromFlags()
	if err != nil {
		panic(fmt.Errorf("failed to read private key: %w", err))
	}
	var chainID *big.Int
	if value, ok := chainIDMap[strings.ToLower(*netnameFlag)]; ok {
		chainID = big.NewInt(int64(value))
	}

	txBuilder, err := chain.NewTxBuilder(*providerFlag, privateKey, chainID)
	if err != nil {
		panic(fmt.Errorf("cannot connect to web3 provider: %w", err))
	}

	port, err := strconv.Atoi(os.Getenv("HTTP_PLATFORM_PORT"))
	if err != nil {
		port = 8080
	}
	httpPortFlag := flag.Int("httpport", port, "Listener port to serve HTTP connection")

	faucetAmount, err := strconv.Atoi(os.Getenv("FAUCET_AMOUNT"))
	if err != nil {
		faucetAmount = 10000
	}
	payoutFlag := flag.Int("faucet.amount", faucetAmount, "Number of Ethers to transfer per user request")

	config := server.NewConfig(*netnameFlag, *symbolFlag, *httpPortFlag, *intervalFlag, *payoutFlag, *proxyCntFlag, *hcaptchaSiteKeyFlag, *hcaptchaSecretFlag)
	go server.NewServer(txBuilder, config).Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func getPrivateKeyFromFlags() (*ecdsa.PrivateKey, error) {
	if *privKeyFlag != "" {
		hexkey := *privKeyFlag
		if chain.Has0xPrefix(hexkey) {
			hexkey = hexkey[2:]
		}
		return crypto.HexToECDSA(hexkey)
	} else if *keyJSONFlag == "" {
		return nil, errors.New("missing private key or keystore")
	}

	keyfile, err := chain.ResolveKeyfilePath(*keyJSONFlag)
	if err != nil {
		return nil, err
	}
	password, err := os.ReadFile(*keyPassFlag)
	if err != nil {
		return nil, err
	}

	return chain.DecryptKeyfile(keyfile, strings.TrimRight(string(password), "\r\n"))
}
