package chain

import (
	"crypto/ecdsa"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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

func ResolveKeyfilePath(keydir string) (string, error) {
	keydir, _ = filepath.Abs(keydir)
	fileInfo, err := os.Stat(keydir)
	if err != nil {
		return "", err
	}
	if !fileInfo.IsDir() {
		return keydir, nil
	}

	files, err := ioutil.ReadDir(keydir)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), "UTC") {
			return filepath.Join(keydir, file.Name()), nil
		}
	}

	return "", errors.New("keystore file not found ")
}
