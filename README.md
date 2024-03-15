# eth-faucet

[![Build](https://img.shields.io/github/actions/workflow/status/chainflag/eth-faucet/build.yml?branch=main)](https://github.com/chainflag/eth-faucet/actions/workflows/build.yml)
[![Release](https://img.shields.io/github/v/release/chainflag/eth-faucet)](https://github.com/chainflag/eth-faucet/releases)
[![Report](https://goreportcard.com/badge/github.com/chainflag/eth-faucet)](https://goreportcard.com/report/github.com/chainflag/eth-faucet)
[![Go](https://img.shields.io/github/go-mod/go-version/chainflag/eth-faucet)](https://go.dev/)
[![License](https://img.shields.io/github/license/chainflag/eth-faucet)](https://github.com/chainflag/eth-faucet/blob/main/LICENSE)

The faucet is a web application with the goal of distributing small amounts of Ether in private and test networks.

## Features

* Configure the funding account using a private key or keystore
* Implement CAPTCHA verification to prevent abuse
* Rate-limit requests by ETH address and IP address to prevent spam
* Prevent X-Forwarded-For spoofing by specifying the number of reverse proxies

## Get started

### Prerequisites

* Go (version 1.17 or later)
* Node.js

### Installation

1. Clone the repository and navigate to the app’s directory

```bash
git clone https://github.com/chainflag/eth-faucet.git
cd eth-faucet
```

2. Bundle frontend with Vite
```bash
go generate
```

3. Build Go project

```bash
go build -o eth-faucet
```

## Usage

**Use a private key**

```bash
./eth-faucet -httpport 8080 -wallet.provider http://localhost:8545 -wallet.privkey privkey
```

**Use a keystore**

```bash
./eth-faucet -httpport 8080 -wallet.provider http://localhost:8545 -wallet.keyjson keystore -wallet.keypass password.txt
```

### Configuration

You can configure the funding account by using environment variables instead of command-line flags:
```bash
export WEB3_PROVIDER=rpc_endpoint
export PRIVATE_KEY=hex_private_key
```

or

```bash
export WEB3_PROVIDER=rpc_endpoint
export KEYSTORE=keystore_path
echo "your_keystore_password" > `pwd`/password.txt
```

Then run the faucet application without the wallet command-line flags:

```bash
./eth-faucet -httpport 8080
```

**Optional Flags**

The following are the available command-line flags(excluding above wallet flags):

| Flag              | Description                                      | Default Value |
|-------------------|--------------------------------------------------|---------------|
| -httpport         | Listener port to serve HTTP connection           | 8080          |
| -proxycount       | Count of reverse proxies in front of the server  | 0             |
| -faucet.amount    | Number of Gwei to transfer per user request      | 1000000000    |
| -faucet.minutes   | Number of minutes to wait between funding rounds | 1440          |
| -faucet.name      | Network name to display on the frontend          | testnet       |
| -faucet.symbol    | Token symbol to display on the frontend          | ETH           |
| -hcaptcha.sitekey | hCaptcha sitekey                                 |               |
| -hcaptcha.secret  | hCaptcha secret                                  |               |
| -frontend.logo    | Logo URL to display on the frontend              | /gatewayfm-logo.svg |
| -frontend.background | Background to display on the frontend         | /background.jpg |

### Docker deployment

```bash
docker run -d -p 8080:8080 -e WEB3_PROVIDER=rpc_endpoint -e PRIVATE_KEY=hex_private_key chainflag/eth-faucet:1.2.0
```

or

```bash
docker run -d -p 8080:8080 -e WEB3_PROVIDER=rpc_endpoint -e KEYSTORE=keystore_path -v `pwd`/keystore:/app/keystore -v `pwd`/password.txt:/app/password.txt chainflag/eth-faucet:1.2.0
```

## License

Distributed under the MIT License. See LICENSE for more information.
