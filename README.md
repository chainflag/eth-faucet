# eth-faucet

## Quick Start
```bash
docker build -t eth-faucet .
docker run -d -p 8080:8080 -v `pwd`/config.yml:/app/config.yml eth-faucet
```
