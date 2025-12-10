package chain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

func EtherToWei(amount float64) *big.Int {
	if amount < 0 {
		return big.NewInt(0)
	}
	oneEther := decimal.NewFromFloat(1e18)
	result := decimal.NewFromFloat(amount).Mul(oneEther)
	wei, ok := new(big.Int).SetString(result.String(), 10)
	if !ok {
		// Return 0 if conversion fails
		return big.NewInt(0)
	}
	return wei
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
