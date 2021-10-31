# eth-faucet
The faucet is a web application with the goal of distributing small amounts of Ether in private and test networks.

## Get started

### Prerequisites

* Go (1.16 or later)
* Node.js

### Installation

1. Clone the repository and navigate to the app’s directory
```bash
git clone https://github.com/chainflag/eth-faucet.git
cd eth-faucet
```

2. Bundle Frontend web with Rollup
```bash
npm run build
```

3. Build Go project 
```bash
go build -o eth-faucet
```

## Usage

1. Set up Web3 Provider and Private Key
```bash
export WEB3_PROVIDER="rpc endpoint"
export PRIVATE_KEY="hex private key"
```

2. Run the eth faucet application
```bash
./eth-faucet -httpport 8080
```

**Optional Flags**

| Flag        | Description                                      | Default Value
| ----------- | ------------------------------------------------ | -------------
| -chainname  | Network name to display on the frontend          | testnet
| -httpport   | Listener port to serve HTTP connection           | 8080
| -interval   | Number of minutes to wait between funding rounds | 1440
| -payout     | Number of Ethers to transfer per user request    | 1
| -proxycount | Count of reverse proxies in front of the server  | 0
| -queuecap   | Maximum transactions waiting to be sent          | 100

### Docker deployment

* Use private key as sender
```bash
docker run -d -p 8080:8080 -e WEB3_PROVIDER="rpc endpoint" -e PRIVATE_KEY="hex private key" chainflag/eth-faucet:1.0.0
```

* Use keystore file as sender
```bash
docker run -d -p 8080:8080 -e WEB3_PROVIDER="rpc endpoint" -e KEYSTORE="keystore path" -v `pwd`/keystore:/app/keystore -v `pwd`/password.txt:/app/password.txt chainflag/eth-faucet:1.0.0
```

### Heroku deployment

```bash
heroku create
heroku buildpacks:add heroku/nodejs
heroku buildpacks:add heroku/go
heroku config:set WEB3_PROVIDER=rpc endpoint
heroku config:set PRIVATE_KEY=hex private key
git push heroku main
heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

## License

Distributed under the MIT License. See LICENSE.txt for more information.
