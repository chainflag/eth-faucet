package cmd

import (
	"crypto/ecdsa"
	"errors"
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

var (
	appVersion = "v1.2.0"
	chainIDMap = map[string]int{"sepolia": 11155111, "holesky": 17000}

	httpPortFlag = flag.Int("httpport", 8080, "Listener port to serve HTTP connection")
	proxyCntFlag = flag.Int("proxycount", 0, "Count of reverse proxies in front of the server")
	versionFlag  = flag.Bool("version", false, "Print version number")

	payoutFlag     = flag.Int64("faucet.amount", 1000000000, "Number of Gwei to transfer per user request")
	intervalFlag   = flag.Int("faucet.minutes", 1440, "Number of minutes to wait between funding rounds")
	netnameFlag    = flag.String("faucet.name", "testnet", "Network name to display on the frontend")
	symbolFlag     = flag.String("faucet.symbol", "ETH", "Token symbol to display on the frontend")
	logoFlag       = flag.String("frontend.logo", "/gatewayfm-logo.svg", "Logo to display on the frontend")
	backgroundFlag = flag.String("frontend.background", "/background.jpg", "Background to display on the frontend")
	keyJSONFlag    = flag.String("wallet.keyjson", os.Getenv("KEYSTORE"), "Keystore file to fund user requests with")
	keyPassFlag    = flag.String("wallet.keypass", "password.txt", "Passphrase text file to decrypt keystore")
	privKeyFlag    = flag.String("wallet.privkey", os.Getenv("PRIVATE_KEY"), "Private key hex to fund user requests with")
	providerFlag   = flag.String("wallet.provider", os.Getenv("WEB3_PROVIDER"), "Endpoint for Ethereum JSON-RPC connection")

	frontendTypeFlag = flag.String("frontend.type", "redesign", "Type of frontend to generate. Values enum: 'base', 'redesign'.")
	paidCustomerFlag = flag.Bool("faucet.paidcustomer", false, "Whether the faucet belongs to the paid customer")

	hcaptchaSiteKeyFlag = flag.String("hcaptcha.sitekey", os.Getenv("HCAPTCHA_SITEKEY"), "hCaptcha sitekey")
	hcaptchaSecretFlag  = flag.String("hcaptcha.secret", os.Getenv("HCAPTCHA_SECRET"), "hCaptcha secret")
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
	config := server.NewConfig(
		*netnameFlag,
		*symbolFlag,
		*httpPortFlag,
		*intervalFlag,
		*proxyCntFlag,
		*payoutFlag,
		*hcaptchaSiteKeyFlag,
		*hcaptchaSecretFlag,
		*logoFlag,
		*backgroundFlag,
		*frontendTypeFlag,
		*paidCustomerFlag,
	)

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
