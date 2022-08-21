package main

import (
	"github.com/chainflag/eth-faucet/cmd"
)

//go:generate npm run build
func main() {
	cmd.Execute()
}
