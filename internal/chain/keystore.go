package chain

import (
	"crypto/ecdsa"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func DecryptPrivateKey(keyfile, password string) (*ecdsa.PrivateKey, error) {
	jsonBytes, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey(jsonBytes, password)
	if err != nil {
		return nil, err
	}

	return key.PrivateKey, nil
}
