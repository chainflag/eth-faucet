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

	payoutFlag   = flag.Float64("faucet.amount", 1, "Number of Ethers to transfer per user request")
	intervalFlag = flag.Int("faucet.minutes", 1440, "Number of minutes to wait between funding rounds")
	netnameFlag  = flag.String("faucet.name", "testnet", "Network name to display on the frontend")
	symbolFlag   = flag.String("faucet.symbol", "ETH", "Token symbol to display on the frontend")

	keyJSONFlag     = flag.String("wallet.keyjson", os.Getenv("KEYSTORE"), "Keystore file to fund user requests with")
	keyPassFlag     = flag.String("wallet.keypass", "password.txt", "Passphrase text file to decrypt keystore")
	privKeyFlag     = flag.String("wallet.privkey", os.Getenv("PRIVATE_KEY"), "Private key hex to fund user requests with")
	providerFlag    = flag.String("wallet.provider", os.Getenv("WEB3_PROVIDER"), "Endpoint for Ethereum JSON-RPC connection")
	loadBalanceFlag = flag.Bool("wallet.loadbalance", false, "Enable load balancing across multiple addresses")

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
	var privateKeys []*ecdsa.PrivateKey
	var err error

	if *privKeyFlag != "" {
		// Handle multiple private keys separated by commas
		hexkeys := strings.Split(*privKeyFlag, ",")
		for _, hexkey := range hexkeys {
			hexkey = strings.TrimSpace(hexkey)
			if chain.Has0xPrefix(hexkey) {
				hexkey = hexkey[2:]
			}
			pk, err := crypto.HexToECDSA(hexkey)
			if err != nil {
				panic(fmt.Errorf("failed to parse private key: %w", err))
			}
			privateKeys = append(privateKeys, pk)
		}
	} else if *keyJSONFlag != "" {
		// Handle multiple keystore files separated by commas
		keystoreFiles := strings.Split(*keyJSONFlag, ",")
		for _, keystoreFile := range keystoreFiles {
			keystoreFile = strings.TrimSpace(keystoreFile)
			keyfile, err := chain.ResolveKeyfilePath(keystoreFile)
			if err != nil {
				panic(fmt.Errorf("failed to resolve keyfile path: %w", err))
			}
			password, err := os.ReadFile(*keyPassFlag)
			if err != nil {
				panic(fmt.Errorf("failed to read password file: %w", err))
			}
			pk, err := chain.DecryptKeyfile(keyfile, strings.TrimRight(string(password), "\r\n"))
			if err != nil {
				panic(fmt.Errorf("failed to decrypt keyfile: %w", err))
			}
			privateKeys = append(privateKeys, pk)
		}
	} else {
		panic(errors.New("missing private key or keystore"))
	}

	if len(privateKeys) == 0 {
		panic(errors.New("no valid private keys provided"))
	}

	var chainID *big.Int
	if value, ok := chainIDMap[strings.ToLower(*netnameFlag)]; ok {
		chainID = big.NewInt(int64(value))
	}

	var txBuilder chain.TxBuilder
	if len(privateKeys) > 1 && *loadBalanceFlag {
		// Use load balancer with multiple keys
		txBuilder, err = chain.NewTxBuilderLoadBalancer(*providerFlag, privateKeys, chainID)
	} else {
		// Use single key (first one if multiple were provided)
		txBuilder, err = chain.NewTxBuilder(*providerFlag, privateKeys[0], chainID)
	}

	if err != nil {
		panic(fmt.Errorf("cannot connect to web3 provider: %w", err))
	}

	config := server.NewConfig(*netnameFlag, *symbolFlag, *httpPortFlag, *intervalFlag, *proxyCntFlag, *payoutFlag, *hcaptchaSiteKeyFlag, *hcaptchaSecretFlag, *loadBalanceFlag)
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
