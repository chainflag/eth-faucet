package chain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

func EtherToWei(amount float64) *big.Int {
	oneEther := decimal.NewFromFloat(1e18)
	result := decimal.NewFromFloat(amount).Mul(oneEther)
	wei, _ := new(big.Int).SetString(result.String(), 10)
	return wei
}

func GweiToWei(amount int64) *big.Int {
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)
	return new(big.Int).Mul(big.NewInt(amount), multiplier)
}

func Has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func IsValidAddress(address string, checksummed bool) bool {
	if !common.IsHexAddress(address) {
		return false
	}
	return !checksummed || common.HexToAddress(address).Hex() == address
}
