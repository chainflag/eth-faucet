# eth-faucet
The faucet is a web application with the goal of distributing small amounts of Ether in private and test networks.

## Get started
* Run faucet by using private key
```bash
docker run -d -p 8080:8080 -e WEB3_PROVIDER="rpc endpoint" -e PRIVATE_KEY="hex private key" chainflag/eth-faucet -apiport 8080
```

* Run faucet by using keystore file
```bash
docker run -d -p 8080:8080 -e WEB3_PROVIDER="rpc endpoint" -e KEYSTORE="keystore path" -e PASSWORD="keystore pass" -v `pwd`/keystore:/app/keystore chainflag/eth-faucet -apiport 8080
```

**Optional Flags**

| flag      | Description                                      | Default Value
| --------- | ------------------------------------------------ | -------------
| -apiport  | Listener port to serve HTTP connection           | 8080
| -interval | Number of minutes to wait between funding rounds | 1440
| -payout   | Number of Ethers to transfer per user request    | 1
| -queuecap | Maximum transactions waiting to be sent          | 100

## Development
### Prerequisites

* Go (1.16 or later)
* Node.js

### Build

1. Clone the repository and navigate to the app’s directory
```bash
git clone https://github.com/chainflag/eth-faucet.git
cd eth-faucet
```

2. Bundle Front-end web with Rollup
```bash
cd web && npm install
npm run build
```
_For more details, please refer to the [web readme](https://github.com/chainflag/eth-faucet/blob/main/web/README.md)_  

3. Build binary application to run
```bash
cd ..
go build -o eth-faucet main.go
export WEB3_PROVIDER=https://ropsten.infura.io
export PRIVATE_KEY=secret
./eth-faucet
```

## License
This project is licensed under the MIT License
