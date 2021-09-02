package chain

import (
	"regexp"

	"github.com/ethereum/go-ethereum/common"
)

func IsValidAddress(address interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := address.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

func ToCheckSumAddress(address interface{}) string {
	switch v := address.(type) {
	case string:
		return common.HexToAddress(v).Hex()
	case common.Address:
		return v.Hex()
	default:
		return ""
	}
}
